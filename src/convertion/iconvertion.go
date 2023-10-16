package convertion

type IConverter interface {
	Start() error
	Convert(*ConvertionRequest) (*ConvertionResponse, error)
	Terminate() error
}

func NewConverter() (IConverter, error) {
	return NewFfmpegConverter()
}
