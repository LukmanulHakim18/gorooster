package models

import "fmt"

type Headers struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MethodType string

const (
	METHOD_POST   = "POST"
	METHOD_GET    = "GET"
	METHOD_PUT    = "PUT"
	METHOD_PATCH  = "PATCH"
	METHOD_DELETE = "DELETE"
)

type JobAPI struct {
	Endpoint string     `json:"endpoint"`
	Headers  []Headers  `json:"headers"`
	Method   MethodType `json:"method"`
	Data     any        `json:"data"`
}

func (ja JobAPI) IsMethodSupport() bool {
	return ja.Method == METHOD_POST || ja.Method == METHOD_GET || ja.Method == METHOD_PUT || ja.Method == METHOD_DELETE || ja.Method == METHOD_PATCH
}

func (m MethodType) ToString() string {
	return string(m)
}

func (ja JobAPI) Validat() error {
	if ok := ja.IsMethodSupport(); !ok {
		return fmt.Errorf("method not support")
	}
	if ja.Endpoint == "" {
		return fmt.Errorf("endpoint is empty")

	}
	return nil
}
