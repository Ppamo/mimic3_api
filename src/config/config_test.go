package config

import (
	"fmt"
	"ppamo/api/common"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	config_complete = `{
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

func Test_Errors(t *testing.T) {
	var (
		p   *common.ProfileOptionsStruct
		e   *common.AudioEffectStruct
		err error
	)
	err = LoadConfig("/")
	assert.NotNil(t, err)
	assert.Equal(t, "read /: is a directory", fmt.Sprintf("%s", err))
	err = LoadConfigFromContent("")
	assert.NotNil(t, err)
	assert.Equal(t, "unexpected end of JSON input", fmt.Sprintf("%s", err))
	err = LoadConfigFromContent("{")
	assert.NotNil(t, err)
	assert.Equal(t, "unexpected end of JSON input", fmt.Sprintf("%s", err))
	err = LoadConfigFromContent("-")
	assert.NotNil(t, err)
	assert.Equal(t, "invalid character ' ' in numeric literal", fmt.Sprintf("%s", err))
	err = LoadConfigFromContent(`{tesste: 1}`)
	assert.NotNil(t, err)
	assert.Equal(t, "invalid character 't' looking for beginning of object key string", fmt.Sprintf("%s", err))
	err = LoadConfigFromContent(`{"tesste": 1}`)
	assert.Nil(t, err)
	assert.NotNil(t, GetConfig())
	err = LoadConfigFromContent("{}")
	assert.Nil(t, err)
	assert.NotNil(t, GetConfig())
	p, err = GetConfig().GetProfileByName("")
	assert.NotNil(t, err)
	assert.Equal(t, "No profiles loaded", fmt.Sprintf("%s", err))
	assert.Nil(t, p)
	e, err = GetConfig().GetEffectByName("")
	assert.NotNil(t, err)
	assert.Equal(t, "No effects loaded", fmt.Sprintf("%s", err))
	assert.Nil(t, e)
}

func Test_Load(t *testing.T) {
	var err error
	err = LoadConfigFromContent(config_complete)
	assert.Nil(t, err)
	assert.NotNil(t, GetConfig())
	assert.Equal(t, "/tmp", GetConfig().TempFolder)
}

func Test_LoadProfiles(t *testing.T) {
	var (
		p   *common.ProfileOptionsStruct
		err error
	)
	err = LoadConfigFromContent(config_complete)
	assert.Nil(t, err)
	assert.NotNil(t, GetConfig())
	assert.Equal(t, 4, len(GetConfig().Profiles))
	p, err = GetConfig().GetProfileByName("en-female")
	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, &common.ProfileOptionsStruct{
		Name: "en-female", Voice: "en_US/cmu-arctic_low", Speaker: "eey",
		LengthScale: 1, NoiseScale: 0.6, NoiseW: 0.3,
	}, p)
	p, err = GetConfig().GetProfileByName("en-male")
	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, &common.ProfileOptionsStruct{
		Name: "en-male", Voice: "en_US/cmu-arctic_low", Speaker: "bdl",
		LengthScale: 1, NoiseScale: 0.61, NoiseW: 0.32,
	}, p)
	p, err = GetConfig().GetProfileByName("es-female")
	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, &common.ProfileOptionsStruct{
		Name: "es-female", Voice: "es_ES/m-ailabs_low", Speaker: "karen_savage",
		LengthScale: 0.98, NoiseScale: 0.83, NoiseW: 0.24,
	}, p)
	p, err = GetConfig().GetProfileByName("es-male")
	assert.Nil(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, &common.ProfileOptionsStruct{
		Name: "es-male", Voice: "es_ES/m-ailabs_low", Speaker: "victor_villarraza",
		LengthScale: 0.9, NoiseScale: 0.85, NoiseW: 0.46,
	}, p)
	p, err = GetConfig().GetProfileByName("test")
	assert.NotNil(t, err)
	assert.Equal(t, "No profile found", fmt.Sprintf("%s", err))
}

func Test_LoadEffects(t *testing.T) {
	var (
		e   *common.AudioEffectStruct
		err error
	)
	LoadConfigFromContent(config_complete)
	assert.NotNil(t, GetConfig())
	assert.Equal(t, 5, len(GetConfig().Effects))
	e, err = GetConfig().GetEffectByName("dot")
	assert.Nil(t, err)
	assert.NotNil(t, e)
	assert.Equal(t, &common.AudioEffectStruct{
		Name: "dot", Path: "/data/effects/dot.wav",
	}, e)
	e, err = GetConfig().GetEffectByName("long-beep")
	assert.Nil(t, err)
	assert.NotNil(t, e)
	assert.Equal(t, &common.AudioEffectStruct{
		Name: "long-beep", Path: "/data/effects/long-beep.wav",
	}, e)
	e, err = GetConfig().GetEffectByName("short-slide")
	assert.Nil(t, err)
	assert.NotNil(t, e)
	assert.Equal(t, &common.AudioEffectStruct{
		Name: "short-slide", Path: "/data/effects/short-slide.wav",
	}, e)
	e, err = GetConfig().GetEffectByName("tirit")
	assert.Nil(t, err)
	assert.NotNil(t, e)
	assert.Equal(t, &common.AudioEffectStruct{
		Name: "tirit", Path: "/data/effects/tirit.wav",
	}, e)
	e, err = GetConfig().GetEffectByName("tirit-2")
	assert.Nil(t, err)
	assert.NotNil(t, e)
	assert.Equal(t, &common.AudioEffectStruct{
		Name: "tirit-2", Path: "/data/effects/tirit2.wav",
	}, e)
	e, err = GetConfig().GetEffectByName("test")
	assert.NotNil(t, err)
	assert.Equal(t, "No effect found", fmt.Sprintf("%s", err))
}

func Test_ParamsToString(t *testing.T) {
	var (
		p   *common.ProfileOptionsStruct
		err error
	)
	LoadConfigFromContent(config_complete)
	assert.NotNil(t, GetConfig())
	p, err = GetConfig().GetProfileByName("en-female")
	assert.Nil(t, err)
	assert.Equal(t, 10, len(p.ToParamsArray()))
	assert.Equal(t,
		"--voice en_US/cmu-arctic_low --speaker eey --length-scale 1.00 --noise-scale 0.60 --noise-w 0.30",
		strings.Join(p.ToParamsArray(), " "))
	p, err = GetConfig().GetProfileByName("en-male")
	assert.Nil(t, err)
	assert.Equal(t, 10, len(p.ToParamsArray()))
	assert.Equal(t,
		"--voice en_US/cmu-arctic_low --speaker bdl --length-scale 1.00 --noise-scale 0.61 --noise-w 0.32",
		strings.Join(p.ToParamsArray(), " "))
	p, err = GetConfig().GetProfileByName("es-female")
	assert.Nil(t, err)
	assert.Equal(t, 10, len(p.ToParamsArray()))
	assert.Equal(t,
		"--voice es_ES/m-ailabs_low --speaker karen_savage --length-scale 0.98 --noise-scale 0.83 --noise-w 0.24",
		strings.Join(p.ToParamsArray(), " "))
	p, err = GetConfig().GetProfileByName("es-male")
	assert.Nil(t, err)
	assert.Equal(t, 10, len(p.ToParamsArray()))
	assert.Equal(t,
		"--voice es_ES/m-ailabs_low --speaker victor_villarraza --length-scale 0.90 --noise-scale 0.85 --noise-w 0.46",
		strings.Join(p.ToParamsArray(), " "))
}
