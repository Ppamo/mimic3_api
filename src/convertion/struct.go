package convertion

import "fmt"

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
	Profile string `json:"profile"`
	Text    string `json:"text"`
}

type ConvertionResponse struct {
	Body *[]byte `json:"body"`
}

func (req *ConvertionRequest) String() string {
	return fmt.Sprintf("%s: %s", req.Profile, req.Text)
}

func (res *ConvertionResponse) String() string {
	return fmt.Sprintf("len: %d", len(*res.Body))
}
