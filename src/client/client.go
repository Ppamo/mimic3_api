package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"ppamo/api/common"
	"ppamo/api/utils"
	"strings"
)

const (
	url  = "http://localhost:8080/convert"
	mime = "application/json"
)

type SourcesStruct struct {
	Name    string                     `json:"name"`
	Sources []common.AudioSourceStruct `json:"sources"`
}

type RequestStruct struct {
	Requests []SourcesStruct `json:"requests"`
}

func main() {
	var (
		requests RequestStruct
		source   SourcesStruct
		bin      []byte
		filePath string
		req      *http.Request
		res      *http.Response
		client   *http.Client
		decoder  *json.Decoder
		f        *os.File
		r        common.ConvertResponse
		err      error
	)
	if len(os.Args) != 2 {
		fmt.Printf("Usage: app <File path>\n")
		return
	}
	filePath = os.Args[1]
	fmt.Printf("> Reading file %s:\n", filePath)
	jsonFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()
	content, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(content, &requests)
	client = &http.Client{}
	for _, source = range requests.Requests {
		fmt.Printf("> Calling: %s\n%+v\n", source.Name, source.Sources)
		bin, err = json.Marshal(source)
		if err != nil {
			fmt.Printf("ERROR: Failed to marshal source:\n%s\n", err)
			return
		}
		req, err = http.NewRequest("POST", url, bytes.NewReader(bin))
		if err != nil {
			fmt.Printf("ERROR: Failed creating request:\n%s\n", err)
			return
		}
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")
		res, err = client.Do(req)
		if err != nil {
			fmt.Printf("ERROR: Failed request:\n%s\n", err)
			return
		}
		fmt.Printf("> Response: %+v\n", res)
		if res.StatusCode != 200 {
			r = common.ConvertResponse{}
			json.NewDecoder(res.Body).Decode(&r)
			fmt.Printf("ERROR: Response status not ok:\n%+v\n", r)
			return
		}
		filePath = utils.GetTimestampedFileName("./",
			fmt.Sprintf("%s.audio.ogg", strings.ReplaceAll(source.Name, " ", "_")))
		fmt.Printf("> Writing new file %s\n", filePath)
		f, err = os.Create(filePath)
		if err != nil {
			fmt.Printf("ERROR: Failed to create output file:\n%s\n", err)
			return
		}
		defer f.Close()
		defer res.Body.Close()
		decoder = json.NewDecoder(res.Body)
		r = common.ConvertResponse{}
		err = decoder.Decode(&r)
		if err != nil {
			fmt.Printf("ERROR: Failed to decode response body:\n%s\n", err)
			return
		}
		err = binary.Write(f, binary.LittleEndian, r.Body)
		if err != nil {
			fmt.Printf("ERROR: Failed to write new file:\n%s\n", err)
			return
		}
	}
}
