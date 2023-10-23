package convertion

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"ppamo/api/common"
	"ppamo/api/config"
	"ppamo/api/utils"
	"regexp"
	"strings"
	"time"
)

const (
	mimic3Bin = "/home/mimic3/app/.venv/bin/mimic3"
	ffmpegBin = "/usr/bin/ffmpeg"
	voicesDir = "/opt/mimic3-server/voices"
)

type ConvertionHandler struct{}

func NewFfmpegConverter() (IConverter, error) {
	c := ConvertionHandler{}
	if err := c.Start(); err != nil {
		log.Printf("> ERROR:\n%v", err)
		return nil, err
	}
	return &c, nil
}

func (c *ConvertionHandler) Start() error {
	log.Printf("ff> Starting converter")
	return nil
}

func (c *ConvertionHandler) getDuration(logs string) (float64, error) {
	var (
		re1      *regexp.Regexp
		re2      *regexp.Regexp
		text     string
		duration time.Duration
		total    time.Duration = 0
		err      error
	)
	log.Printf(">> Logs:\n%s", logs)
	re1 = regexp.MustCompile(`Duration: [0-9:\.]+,`)
	re2 = regexp.MustCompile(`[0-9][^,]*`)
	for _, line := range re1.FindAllStringSubmatch(logs, -1) {
		log.Printf("1>>.")
		text = re2.FindString(line[0])
		text = strings.Replace(text, ":", "h", 1)
		text = strings.Replace(text, ":", "m", 1)
		text = strings.Replace(text, ".", "s", 1)
		text = fmt.Sprintf("%s0ms", text)
		duration, err = time.ParseDuration(text)
		if err != nil {
			return 0.0, err
		}
		log.Printf("1>> Duration: %s", duration)
		total = total + duration
	}
	log.Printf("2>> Duration: %s", total)
	return total.Seconds(), nil
}

func (c *ConvertionHandler) convertWavToOgg(files []string) (common.ConvertResponse, error) {
	var (
		i        int
		cmd      *exec.Cmd
		path     string
		params   []string
		fcomplex string
		output   []byte
		duration float64
		err      error
	)

	params = []string{"-y"}
	for i, path = range files {
		params = append(params, "-i", path)
		fcomplex = fmt.Sprintf("%s[%d:a]", fcomplex, i)
	}
	path = utils.GetTimestampedFileName(config.GetConfig().TempFolder, "output.ogg")
	fcomplex = fmt.Sprintf("%sconcat=n=%d:v=0:a=1", fcomplex, len(files))
	params = append(params, "-filter_complex", fcomplex, "-c:a", "libopus", "-b:a", "48000", path)
	log.Printf("ff> Executing ffmpeg %s", strings.Join(params, " "))
	cmd = exec.Command(ffmpegBin, params...)
	if output, err = cmd.CombinedOutput(); err != nil {
		log.Printf("ERROR: Failed to start ffmpeg:\n%s\n%v", output, err)
		return common.ConvertResponse{}, err
	}
	log.Printf("cf> ffmpeg output: %s", output)
	duration, err = c.getDuration(string(output))
	defer os.Remove(path)
	if output, err = os.ReadFile(path); err != nil {
		log.Printf("ERROR: Failed to start ffmpeg:\n%s\n%v", output, err)
		return common.ConvertResponse{}, err
	}
	if err != nil {
		log.Printf("ERROR: Failed to calculate duration:\n%v", err)
		return common.ConvertResponse{}, err
	}
	return common.ConvertResponse{Body: output, Duration: duration}, nil
}

func (c *ConvertionHandler) createTTSWav(p *common.ProfileOptionsStruct, text string) (string, error) {
	var (
		params []string
		cmd    *exec.Cmd
		f      *os.File
		stderr bytes.Buffer
		err    error
	)
	params = p.ToParamsArray()
	params = append(params, "--voices-dir", voicesDir, text)
	cmd = exec.Command(mimic3Bin, params...)
	f, err = os.CreateTemp(config.GetConfig().TempFolder, "out.*.wav")
	if err != nil {
		log.Printf("ERROR: Failed to create output temp file:\n%v", err)
		return "", err
	}
	defer f.Close()
	cmd.Stdout = f
	cmd.Stderr = &stderr
	log.Printf("cf> Executing:\n%s %s", mimic3Bin, strings.Join(params, " "))
	log.Printf("cf> Writting output temp file %s", f.Name())
	if err = cmd.Run(); err != nil {
		log.Printf("ERROR: Failed to start mimic3:\n%s\n%v", stderr.String(), err)
		return "", err
	}
	log.Printf("cf> mimic3 Output:\n%s", stderr.String())
	return f.Name(), nil
}

func (c *ConvertionHandler) Convert(req *common.ConvertRequest) (*common.ConvertResponse, error) {
	var (
		source  common.AudioSourceStruct
		wavFile string
		files   []string
		res     common.ConvertResponse
		err     error
	)
	if len(req.Sources) == 0 {
		log.Printf("cf> No sources found!")
		return nil, nil
	}
	log.Printf("ff> Sources:\n%+v\n", req.Sources)
	for _, source = range req.Sources {
		switch common.StringToSourceType(source.Type) {
		case common.SourceText:
			log.Printf("ff> Converting\n%v", source.Text)
			if wavFile, err = c.createTTSWav(source.Profile, source.Text); err != nil {
				log.Printf("ERROR: Failed to create TTS\n%s", err)
				return nil, err
			}
			defer os.Remove(wavFile)
			files = append(files, wavFile)
		case common.SourceEffect:
			files = append(files, source.EffectPath)
		}
	}
	if res, err = c.convertWavToOgg(files); err != nil {
		log.Printf("ERROR: Failed to create Ogg\n%s", err)
		return nil, err
	}
	return &res, nil
}

func (c *ConvertionHandler) Terminate() error {
	log.Printf("ff> Terminating converter")
	return nil
}
