package convertion

import "ppamo/api/common"

type IConverter interface {
	Start() error
	Convert(*common.ConvertRequest) (*common.ConvertResponse, error)
	Terminate() error
}

func NewConverter() (IConverter, error) {
	return NewFfmpegConverter()
}
