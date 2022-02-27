package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var (
	ErrorInvalidQuery = trace.New("ERROR_INVALID_QUERY")
	ErrorInvalidBody  = trace.New("ERROR_INVALID_BODY")
	ErrorParseBody    = trace.New("ERROR_PARSE_BODY")
	ErrorCreateUser   = trace.New("ERROR_CREATE_USER")
	ErrorNotFound     = trace.New("ERROR_NOT_FOUND")
)

func MustMarshal(resp Response) string {
	serialised, _ := json.Marshal(resp)
	return string(serialised)
}

func writeErrorResponse(rw http.ResponseWriter, code int, message string, err *trace.Error) {
	logrus.WithField("error", err.Tracef(message)).Error("request failed")
	rw.WriteHeader(code)
	fmt.Fprint(rw, MustMarshal(RespondWithError(err.Name, message)))
}

func writeSuccessResponse(rw http.ResponseWriter, data interface{}) {
	logrus.WithField("data", data).Debug("request succeeded")
	rw.WriteHeader(http.StatusOK)
	fmt.Fprint(rw, MustMarshal(RespondWithSuccess(data)))
}

func RespondWithError(code, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}

func RespondWithSuccess(data interface{}) Response {
	return Response{
		Code:    "Success",
		Message: "OK",
		Data:    data,
	}
}
