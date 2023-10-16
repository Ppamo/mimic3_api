package convertion

import (
	"fmt"
	"log"
	"ppamo/api/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	configuration = `{
	"temp_folder": "/tmp",
	"profiles":[
		{ "name": "en-female", "voice": "en_US/cmu-arctic_low", "speaker": "eey",
			"length_scale": 1, "noise_scale": 0.6, "noise_w": 0.3 },
		{ "name": "en-male", "voice": "en_US/cmu-arctic_low", "speaker": "bdl",
			"length_scale": 1, "noise_scale": 0.61, "noise_w": 0.32 },
		{ "name": "es-female", "voice": "es_ES/m-ailabs_low", "speaker": "karen_savage",
			"length_scale": 0.98, "noise_scale": 0.83, "noise_w": 0.24 },
		{ "name": "es-male", "voice": "es_ES/m-ailabs_low", "speaker": "victor_villarraza",
			"length_scale": 0.9, "noise_scale": 0.85, "noise_w": 0.46 }
	],
	"effects":[
		{"name": "dot", "path": "/data/effects/dot.wav"},
		{"name": "long-beep", "path": "/data/effects/long-beep.wav"},
		{"name": "short-slide", "path": "/data/effects/short-slide.wav"},
		{"name": "tirit", "path": "/data/effects/tirit.wav"},
		{"name": "tirit-2", "path": "/data/effects/tirit2.wav"}
	]
}`
)

var (
	err       error
	converter IConverter
)

func TestMain(m *testing.M) {
	converter, err = NewConverter()
	if err != nil {
		log.Fatalf("ct> ERROR: Unable to create converter:\n%s", err)
	}
	m.Run()
}

func Test_ConvertError(t *testing.T) {
	var (
		req *ConvertionRequest
		res *ConvertionResponse
		err error
	)
	req = &ConvertionRequest{
		Profile: "",
		Text:    "",
	}
	res, err = converter.Convert(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, "No profiles loaded", fmt.Sprintf("%s", err))
}

func assertRequest(t *testing.T, req *ConvertionRequest) {
	var (
		res      *ConvertionResponse
		filePath string
		err      error
	)
	res, err = converter.Convert(req)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.True(t, len(*res.Body) > 0)
	filePath = fmt.Sprintf("/data/tmp/%s.ogg", req.Profile)
	// err = os.WriteFile(filePath, *res.Body, 0666)
	// assert.Nil(t, err)
}

func Test_Convert(t *testing.T) {
	var (
		req *ConvertionRequest
	)
	config.LoadConfigFromContent(configuration)
	req = &ConvertionRequest{Profile: "es-female", Text: "Hola mi mundo!, como esta usted?"}
	assertRequest(t, req)
	req = &ConvertionRequest{Profile: "es-male", Text: "Hola mi mundo!, como esta usted?"}
	assertRequest(t, req)
	req = &ConvertionRequest{Profile: "en-female", Text: "Hello my world!, how are you?"}
	assertRequest(t, req)
	req = &ConvertionRequest{Profile: "en-male", Text: "Hello my world!, how are you?"}
	assertRequest(t, req)

}
