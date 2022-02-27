package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/sirupsen/logrus"
)

type response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var (
	ErrInvalidQuery = trace.New("ErrInvalidQuery")
	ErrInvalidBody  = trace.New("ErrInvalidBody")
	ErrParseBody    = trace.New("ErrParseBody")
	ErrCreateUser   = trace.New("ErrCreateUser")
	ErrNotFound     = trace.New("ErrNotFound")
)

func serializeResponse(resp response) string {
	serialised, _ := json.Marshal(resp)
	return string(serialised)
}

func writeErrorResponse(rw http.ResponseWriter, code int, message string, errWrapper *trace.Error, err error) {
	var traced *trace.Error
	if err != nil {
		traced = errWrapper.Trace(err)
	} else {
		traced = errWrapper.Tracef(message)
	}

	logrus.WithFields(logrus.Fields{"error": traced}).Error(message)
	rw.WriteHeader(code)
	fmt.Fprint(rw, RespondWithError(errWrapper.Name, message))
}

func writeSuccessResponse(rw http.ResponseWriter, data interface{}) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprint(rw, RespondWithSuccess(data))
}

func RespondWithError(code, message string) string {
	return serializeResponse(response{
		Code:    code,
		Message: message,
	})
}

func RespondWithSuccess(data interface{}) string {
	return serializeResponse(response{
		Code:    "Success",
		Message: "OK",
		Data:    data,
	})
}
