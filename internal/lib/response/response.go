package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	statusOk  = "OK"
	statusErr = "ERROR"
)

func Ok() Response {
	return Response{
		Status: statusOk,
	}
}

func Error(msg string) Response {
	return Response{
		Status: statusErr,
		Error:  msg,
	}
}

func ValidationError(err validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range err {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("%s is a required field", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("%s is not a valid URL", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("%s is not a valid", err.Field()))
		}
	}

	return Response{
		Status: statusErr,
		Error:  strings.Join(errMsgs, ", "),
	}
}
