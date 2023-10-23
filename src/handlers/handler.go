package handlers

import (
	"fmt"
	"log"
	"ppamo/api/common"
	"ppamo/api/config"
	"ppamo/api/convertion"
	"strings"
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
		res    *common.ConvertResponse
		effect *common.AudioEffectStruct
		i      int
		err    error
	)
	for i = range req.Sources {
		switch common.StringToSourceType(req.Sources[i].Type) {
		case common.SourceText:
			if len(strings.TrimSpace(req.Sources[i].Text)) == 0 {
				return &common.ConvertResponse{Status: 401, Description: "Bad Request: Text empty"},
					fmt.Errorf("Content empty: %s", err)
			}
			req.Sources[i].Profile, err = config.GetConfig().GetProfileByName(req.Sources[i].ProfileName)
			if err != nil {
				return &common.ConvertResponse{Status: 401, Description: "Bad Request: Invalid profile"},
					fmt.Errorf("Profile not found: %s", err)
			}
		case common.SourceEffect:
			effect, err = config.GetConfig().GetEffectByName(req.Sources[i].EffectName)
			if err != nil {
				return &common.ConvertResponse{Status: 401, Description: "Bad Request: Invalid effect"},
					fmt.Errorf("Effect not found: %s", err)
			}
			req.Sources[i].EffectPath = effect.Path
		default:
			return &common.ConvertResponse{Status: 401, Description: "Bad Request: Invalid Source Type"},
				fmt.Errorf("Unknown Source Type: %s", err)
		}
	}
	if res, err = h.converter.Convert(req); err != nil {
		log.Printf("hl> ERROR: Failed to create audio:\n%s", err)
		return &common.ConvertResponse{Status: 401, Description: "Bad Request: Failed to create audio"},
			fmt.Errorf("Failed to create audio: %s", err)
	}
	res.Status = 200
	return res, nil
}
