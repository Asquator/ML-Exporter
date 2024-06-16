package encode

import (
	"encoding/json"
	"errors"

	ort "github.com/yalue/onnxruntime_go"
)

var (
	ErrInvalidEncoding = errors.New("invalid encoding")
)

type InputMap map[string]InputEntry
type InputUnmarshaler func([]byte) (InputValue, error)

var unmarshalers = map[string]InputUnmarshaler{
	"tensor": func(b []byte) (InputValue, error) {
		var t Tensor
		err := json.Unmarshal(b, &t)

		if err != nil {
			return nil, ErrInvalidEncoding
		}

		return &t, nil
	},
}

type InputValue interface {
	ToONNXValue() (ort.ArbitraryTensor, error)
}

type InputEntry struct {
	Type  string
	Value InputValue
}

func (iv *InputEntry) UnmarshalJSON(b []byte) error {
	type genericInputEntry struct {
		Type  string          `json:"type"`
		Value json.RawMessage `json:"value"`
	}

	var entry genericInputEntry
	err := json.Unmarshal(b, &entry)

	if err != nil {
		return ErrInvalidEncoding
	}

	unmarshaler, ok := unmarshalers[entry.Type]

	if !ok {
		return ErrInvalidEncoding
	}

	iv.Value, err = unmarshaler(entry.Value)

	if err != nil {
		return ErrInvalidEncoding
	}

	return nil
}
