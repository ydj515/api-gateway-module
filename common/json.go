package common

import (
	"log"

	"github.com/bytedance/sonic" // 빠른 marshal, unmarshal을 위해 사용
)

type jsonHandler struct {
	marshal   func(v interface{}) ([]byte, error)
	unmarshal func(data []byte, v interface{}) error
}

var JsonHandler jsonHandler

func init() {
	JsonHandler = jsonHandler{
		marshal:   sonic.Marshal,
		unmarshal: sonic.Unmarshal,
	}
}

// Marshal 에러 로깅을 위해 jsonHandler를 사용
func (j jsonHandler) Marshal(v interface{}) ([]byte, error) {
	bytes, err := j.marshal(v)

	if err != nil {
		log.Println("Failed to marshal", "err", err.Error())
		return nil, err
	}
	return bytes, nil
}

// Unmarshal 에러 로깅을 위해 jsonHandler를 사용
func (j jsonHandler) Unmarshal(data []byte, v interface{}) error {
	err := j.unmarshal(data, v)

	if err != nil {
		log.Println("Failed to unmarshal", "err", err.Error())
		return err
	}
	return nil
}
