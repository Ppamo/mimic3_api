package convertion

import (
	"fmt"
	"ppamo/api/common"
)

type SpeakerProfile int32

const (
	ProfileEnglishFemale SpeakerProfile = iota
	ProfileEnglishMale
	ProfileSpanishFemale
	ProfileSpanishMale
)

func (p SpeakerProfile) String() string {
	return [...]string{"English-Female", "English-Male", "Spanish-Female", "Spanish-Male"}[p]
}

type ConvertionRequest struct {
	Sources []common.AudioSourceStruct `json:"sources"`
}

type ConvertionResponse struct {
	Body *[]byte `json:"body"`
}

func (req *ConvertionRequest) String() string {
	return fmt.Sprintf("%+v", req)
}

func (res *ConvertionResponse) String() string {
	return fmt.Sprintf("len: %d", len(*res.Body))
}
