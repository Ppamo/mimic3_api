package common

import (
	"fmt"
	"strings"
)

type SourceType int32

const (
	SourceUnknown SourceType = iota
	SourceText
	SourceEffect
)

func (st SourceType) String() string {
	return [...]string{"Unknown", "Text", "Effect"}[st]
}

func StringToSourceType(name string) SourceType {
	values := map[string]SourceType{
		"unknown": SourceUnknown,
		"text":    SourceText,
		"effect":  SourceEffect,
	}
	if element, ok := values[strings.ToLower(name)]; ok {
		return element
	}
	return SourceUnknown
}

type AudioSourceStruct struct {
	Type        string                `json:"type"`
	ProfileName string                `json:"profile,omitempty"`
	EffectName  string                `json:"effect,omitempty"`
	Text        string                `json:"text"`
	Profile     *ProfileOptionsStruct `json:"profile_details,omitempty"`
	EffectPath  string                `json:"filepath,omitempty"`
}

type ConvertRequest struct {
	Sources []AudioSourceStruct `json:"sources"`
}

type ConvertResponse struct {
	Status      int     `json:"status"`
	Description string  `json:"description,omitempty"`
	Body        []byte  `json:"body,omiempty"`
	Duration    float64 `json:"duration"`
}

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

type ProfilesResponse struct {
	Status      int                    `json:"status"`
	Description string                 `json:"description,omitempty"`
	Profiles    []ProfileOptionsStruct `json:"profiles"`
}

type EffectsResponse struct {
	Status      int                 `json:"status"`
	Description string              `json:"description,omitempty"`
	Effects     []AudioEffectStruct `json:"effects"`
}

type DefaultError struct {
	Status      int    `json:"status"`
	Description string `json:"description"`
}
