package handlers

import (
	"log"
	"ppamo/api/common"
	"ppamo/api/config"
	"ppamo/api/convertion"
)

type HandlerStruct struct {
	converter convertion.IConverter
}

func NewBasicHandler(c *convertion.IConverter) (IHandler, error) {
	var (
		h HandlerStruct
	)
	h = HandlerStruct{converter: *c}
	return &h, nil
}

func (h *HandlerStruct) GetProfiles() (*common.ProfilesResponse, error) {
	return &common.ProfilesResponse{
		Status:      200,
		Description: "OK",
		Profiles:    config.GetConfig().Profiles,
	}, nil
}

func (h *HandlerStruct) GetEffects() (*common.EffectsResponse, error) {
	return &common.EffectsResponse{
		Status:      200,
		Description: "OK",
		Effects:     config.GetConfig().Effects,
	}, nil
}

func (h *HandlerStruct) Convert(req *common.ConvertRequest) (*common.ConvertResponse, error) {
	var (
		res *common.ConvertResponse
		err error
	)
	if res, err = h.converter.Convert(req); err != nil {
		log.Printf("hl> ERROR: Failed to create audio:\n%s", err)
		return nil, err
	}
	return res, nil
}
