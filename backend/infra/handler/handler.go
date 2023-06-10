package handler

import (
	"code-shooting/infra/common"
	"net/http"
)

const (
	SUCCESS = "success"
	FAILURE = "failure"
)

func Failure(cause string) *common.Response {
	return &common.Response{
		Result: FAILURE,
		Detail: cause,
	}
}

func FailureNotFound(cause string) *common.Response {
	return &common.Response{
		Result: FAILURE,
		Code:   http.StatusNotFound,
		Detail: cause,
	}
}

func Success(any interface{}) *common.Response {
	return &common.Response{
		Result: SUCCESS,
		Detail: any,
	}
}

func SuccessGet(any interface{}) *common.Response {
	return &common.Response{
		Result: SUCCESS,
		Code:   http.StatusOK,
		Detail: any,
	}
}

func SuccessCreated(any interface{}) *common.Response {
	return &common.Response{
		Result: SUCCESS,
		Code:   http.StatusCreated,
		Detail: any,
	}
}
