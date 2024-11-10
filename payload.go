package asynqplus

import (
	"encoding/json"
	"strconv"
)

const (
	_paramName  = "param"
	_resultName = "result"
)

type Payload struct {
	Params map[string]json.RawMessage `json:"params"`
}

func paramName(i int) string {
	return _paramName + strconv.Itoa(i)
}

type Result struct {
	Result map[string]json.RawMessage `json:"result"`
}

func resultName(i int) string {
	return _resultName + strconv.Itoa(i)
}
