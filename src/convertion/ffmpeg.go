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
	"strings"
)

const (
	mimic3Bin = "/home/mimic3/app/.venv/bin/mimic3"
	ffmpegBin = "/usr/bin/ffmpeg"
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

func (c *ConvertionHandler) convertWavToOgg(files []string) ([]byte, error) {
	var (
		i        int
		cmd      *exec.Cmd
		path     string
		params   []string
		fcomplex string
		output   []byte
		err      error
	)

	params = []string{"-y"}
	for i, path = range files {
		params = append(params, "-i", path)
		fcomplex = fmt.Sprintf("%s[%d:a]", fcomplex, i)
	}
	path = utils.GetTimestampedFileName(config.GetConfig().TempFolder, "output.ogg")
	fcomplex = fmt.Sprintf("%sconcat=n=%d:v=0:a=1", fcomplex, len(files))
	params = append(params, "-filter_complex", fcomplex, "-c:a", "libvorbis", path)
	log.Printf("ff> Executing ffmpeg %s", strings.Join(params, " "))
	cmd = exec.Command(ffmpegBin, params...)
	if output, err = cmd.CombinedOutput(); err != nil {
		log.Printf("ERROR: Failed to start ffmpeg:\n%s\n%v", output, err)
		return nil, err
	}
	log.Printf("cf> ffmpeg output: %s", output)
	defer os.Remove(path)
	if output, err = os.ReadFile(path); err != nil {
		log.Printf("ERROR: Failed to start ffmpeg:\n%s\n%v", output, err)
		return nil, err
	}
	return output, nil
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
	params = append(params, text)
	cmd = exec.Command(mimic3Bin, params...)
	f, err = os.CreateTemp(config.GetConfig().TempFolder, "out.*.wav")
	if err != nil {
		log.Printf("ERROR: Failed to create output temp file:\n%v", err)
		return "", err
	}
	defer f.Close()
	cmd.Stdout = f
	cmd.Stderr = &stderr
	log.Printf("cf> Writting output temp file %s", f.Name())
	if err = cmd.Run(); err != nil {
		log.Printf("ERROR: Failed to start mimic3:\n%s\n%v", stderr.String(), err)
		return "", err
	}
	return f.Name(), nil
}

func (c *ConvertionHandler) Convert(req *common.ConvertRequest) (*common.ConvertResponse, error) {
	/*
		type AudioSourceStruct struct {
			Source  SourceType `json:"source"`
			Profile string     `json:"profile,omitempty"`
			Content string     `json:"value"`
		}
		type ConvertRequest struct {
			Sources []AudioSourceStruct `json:"sources"`
		}

		type ConvertResponse struct {
			Status      int    `json:"status"`
			Description string `json:"description,omitempty"`
			Body        []byte `json:"body"`
		}
	*/
	var (
		source  common.AudioSourceStruct
		wavFile string
		files   []string
		output  []byte
		p       *common.ProfileOptionsStruct
		e       *common.AudioEffectStruct
		res     common.ConvertResponse
		err     error
	)
	if len(req.Sources) == 0 {
		log.Printf("cf> No sources found!")
		return nil, nil
	}
	log.Printf("ff> Sources:\n%+v\n", req.Sources)
	for _, source = range req.Sources {
		switch common.StringToSourceType(source.Source) {
		case common.SourceText:
			if p, err = config.GetConfig().GetProfileByName(source.Profile); err != nil {
				log.Printf("ERROR: Failed to get profile '%s'\n%s", source.Profile, err)
				return nil, err
			}
			log.Printf("ff> Converting\n%v", source.Content)
			if wavFile, err = c.createTTSWav(p, source.Content); err != nil {
				log.Printf("ERROR: Failed to create TTS\n%s", err)
				return nil, err
			}
			defer os.Remove(wavFile)
			files = append(files, wavFile)
		case common.SourceEffect:
			e, err = config.GetConfig().GetEffectByName(source.Content)
			if err != nil {
				log.Printf("ERROR: Failed to load '%s' effect", source.Content, err)
				return nil, err
			}
			files = append(files, e.Path)
		}
	}
	if output, err = c.convertWavToOgg(files); err != nil {
		log.Printf("ERROR: Failed to create Ogg\n%s", err)
		return nil, err
	}
	res = common.ConvertResponse{Body: output}
	return &res, nil
}

func (c *ConvertionHandler) Terminate() error {
	log.Printf("ff> Terminating converter")
	return nil
}
