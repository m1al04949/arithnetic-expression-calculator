package response

import (
	"net/http"
)

type Response struct {
	Status     string `json:"status"`
	HttpStatus int    `json:"httpstatus"`
	ID         int    `json:"id"`
	Token      string `json:"token"`
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

func Authorization(id int, token string) Response {
	return Response{
		Status:     StatusOK,
		HttpStatus: http.StatusOK,
		ID:         id,
		Token:      token,
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

func ErrorAuthorization(msg string) Response {
	return Response{
		Status:     StatusError,
		HttpStatus: http.StatusUnauthorized,
		Error:      msg,
	}
}

func ErrorRegistration(msg string) Response {
	return Response{
		Status:     StatusError,
		HttpStatus: http.StatusConflict,
		Error:      msg,
	}
}
