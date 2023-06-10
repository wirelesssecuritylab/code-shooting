package errcode

import (
	"net/http"

	"github.com/pkg/errors"
)

type codeError struct {
	errCode    string
	message    string
	statusCode int
}

func (s *codeError) Error() string {
	return s.message
}

func newCodeErr(e string, m string, s int) error {
	return &codeError{errCode: e, message: m, statusCode: s}
}

var (
	ErrParamError      = newCodeErr("1001", "parameter error", http.StatusBadRequest)
	ErrFileSystemError = newCodeErr("1002", "file system access error", http.StatusInternalServerError)
	ErrDBAccessError   = newCodeErr("1003", "data base access error", http.StatusInternalServerError)
	ErrRecordNotFound  = newCodeErr("1004", "record not found", http.StatusNotFound)
	ErrUnknownError    = newCodeErr("1099", "unknown error", http.StatusInternalServerError)

	ErrUserNotExist     = newCodeErr("1101", "user not exist", http.StatusNotFound)
	ErrUnauthorized     = newCodeErr("1102", "unauthorized", http.StatusUnauthorized)
	ErrPermissionDenied = newCodeErr("1103", "permission denied", http.StatusForbidden)

	ErrInvalidAnswerFile    = newCodeErr("1201", "invalid answer file", http.StatusBadRequest)
	ErrStdAnswerNotExist    = newCodeErr("1202", "standard answer not exist", http.StatusServiceUnavailable)
	ErrScoreAnswerFileError = newCodeErr("1203", "score answer file error", http.StatusBadRequest)
)

type RspMsg struct {
	ErrCode string `json:"errCode,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewRspMsg(c, m string) *RspMsg {
	return &RspMsg{ErrCode: c, Message: m}
}

func SuccMsg() *RspMsg {
	return NewRspMsg("", "success")
}

func ToErrRsp(err error) (int, *RspMsg) {
	if cerr, ok := errors.Cause(err).(*codeError); ok {
		return cerr.statusCode, NewRspMsg(cerr.errCode, err.Error())
	}
	return http.StatusInternalServerError, NewRspMsg("1099", err.Error())
}

func SameCause(err1, err2 error) bool {
	return errors.Cause(err1) == errors.Cause(err2)
}
