package common

import "fmt"

type AudioEffectStruct struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type ProfileOptionsStruct struct {
	Name        string  `json:"name"`
	Voice       string  `json:"voice"`
	Speaker     string  `json:"speaker"`
	LengthScale float32 `json:"length_scale"`
	NoiseScale  float32 `json:"noise_scale"`
	NoiseW      float32 `json:"noise_w"`
}

func (p *ProfileOptionsStruct) ToParamsArray() []string {
	return []string{"--voice", p.Voice, "--speaker", p.Speaker,
		"--length-scale", fmt.Sprintf("%.2f", p.LengthScale),
		"--noise-scale", fmt.Sprintf("%.2f", p.NoiseScale),
		"--noise-w", fmt.Sprintf("%.2f", p.NoiseW)}
}
