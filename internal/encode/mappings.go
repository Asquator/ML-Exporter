package encode

import (
	"errors"

	ort "github.com/yalue/onnxruntime_go"
)

var (
	ErrMappingFailed = errors.New("could not map input to ONNX")
)

func NewTensorFromInput(t Tensor) (*ort.Tensor[float32], error) {
	tensor, err := ort.NewTensor(ort.NewShape(t.Shape...), t.Data)

	if err != nil {
		return nil, errors.Join(err, ErrMappingFailed)
	}

	return tensor, nil
}
