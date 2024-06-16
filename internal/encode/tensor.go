package encode

import (
	"errors"
	"fmt"

	ort "github.com/yalue/onnxruntime_go"
)

type numeric interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

type Tensor struct {
	Shape []int64   `json:"shape"`
	Data  []float32 `json:"data"`
}

func (t *Tensor) ToONNXValue() (ort.ArbitraryTensor, error) {
	tensor, err := ort.NewTensor(ort.NewShape(t.Shape...), t.Data)

	if err != nil {
		return nil, errors.Join(err, ErrInvalidEncoding)
	}

	return ort.ArbitraryTensor(tensor), nil
}

func (t Tensor) String() string {
	return fmt.Sprintf("shape: %v, data: %v", t.Shape, t.Data)
}
