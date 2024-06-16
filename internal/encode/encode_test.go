package encode

import (
	"encoding/json"
	"fmt"
	"testing"

	ort "github.com/yalue/onnxruntime_go"
)

func TestMain(m *testing.M) {
	ort.SetSharedLibraryPath("/usr/lib/libonnxruntime.so")

	ort.InitializeEnvironment()

	m.Run()
}

func Test_decode(t *testing.T) {
	var m InputMap

	jsonString := `
	{
		"inp1" : {
			"type" : "tensor",
			"value" : {
				"shape" : [1, 4],
				"data" : [0.3,4,2,43,3]
			}
		},

		"inp2" : {
			"type" : "tensor",
			"value" : {
				"shape" : [1, 5],
				"data" :  [1,1,1,1,1]
			}
		}
	}`

	err := json.Unmarshal([]byte(jsonString), &m)

	if err != nil {
		t.Error(err)
	}

	for k, v := range m {
		fmt.Println(k, v)

		tensor, err := v.Value.ToONNXValue()
		if err != nil {
			t.Error("cannot convert to ONNX tensor", err)
		}

		tens := tensor.(*ort.Tensor[float32])

		fmt.Println(tens, tens.GetData())

	}
}
