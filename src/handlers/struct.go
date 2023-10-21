package handlers

import (
	"ppamo/api/common"
)

type DefaultResponse struct {
	Status      int    `json:"status"`
	Description string `json:"description,omitempty"`
}

type ConvertRequest struct {
	Sources []common.AudioSourceStruct `json:"sources"`
}

type ConvertResponse struct {
	Status      int    `json:"status"`
	Description string `json:"description,omitempty"`
	Body        []byte `json:"body"`
}

func (req *ConvertRequest) ToConvertionRequest() *common.ConvertRequest {
	return &common.ConvertRequest{Sources: req.Sources}
}
