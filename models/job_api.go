package models

import (
	"fmt"
	"strings"
)

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

// Check if method support or not
func (ja JobAPI) IsMethodSupport() bool {
	return ja.Method == METHOD_POST || ja.Method == METHOD_GET || ja.Method == METHOD_PUT || ja.Method == METHOD_DELETE || ja.Method == METHOD_PATCH
}

// parsing  method to string
func (m MethodType) ToString() string {
	return string(m)
}

// Validate data Job API
func (ja *JobAPI) Validate() error {
	methodType := strings.ToUpper(ja.Method.ToString())
	ja.Method = MethodType(methodType)
	if ok := ja.IsMethodSupport(); !ok {
		return fmt.Errorf("method not support")
	}
	if ja.Endpoint == "" {
		return fmt.Errorf("endpoint is empty")

	}
	return nil
}
