package handlers

import (
	"ppamo/api/common"
	"ppamo/api/convertion"
)

type IHandler interface {
	Convert(req *common.ConvertRequest) (*common.ConvertResponse, error)
	GetProfiles() (*common.ProfilesResponse, error)
	GetEffects() (*common.EffectsResponse, error)
}

func NewHandler(converter *convertion.IConverter) (IHandler, error) {
	return NewBasicHandler(converter)
}
