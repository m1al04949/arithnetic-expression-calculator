package response

import (
	"net/http"
)

type Response struct {
	Status     string `json:"status"`
	HttpStatus int    `json:"httpstatus"`
	Error      string `json:"error,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status:     StatusOK,
		HttpStatus: http.StatusOK,
	}
}

func ErrorRequest(msg string) Response {
	return Response{
		Status:     StatusError,
		HttpStatus: http.StatusBadRequest,
		Error:      msg,
	}
}

func ErrorExpression(msg string) Response {
	return Response{
		Status:     StatusError,
		HttpStatus: http.StatusBadRequest,
		Error:      msg,
	}
}

func ErrorServer(msg string) Response {
	return Response{
		Status:     StatusError,
		HttpStatus: http.StatusInternalServerError,
		Error:      msg,
	}
}
