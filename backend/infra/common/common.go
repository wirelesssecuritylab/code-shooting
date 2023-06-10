package common

import "errors"

type Response struct {
	Result string      `json:"result"`
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Detail interface{} `json:"detail"`
}

var ErrInconsistent = errors.New("get or query result from HR is inconsistent with the request")

const (
	ContentType = "application/json"
)
