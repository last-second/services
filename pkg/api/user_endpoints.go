package api

import (
	"io"
	"net/http"

	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/last-second/services/pkg/config"
	"github.com/last-second/services/pkg/user"
	"github.com/sirupsen/logrus"
)

func createUser(rw http.ResponseWriter, r *http.Request) {
	if r.ContentLength < 1 {
		writeErrorResponse(rw, http.StatusBadRequest, "Body cannot be empty", ErrorInvalidBody)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(rw, http.StatusInternalServerError, "Body could not be read", trace.Guarantee(err).Add("body", body))
		return
	}

	partialUser, err := user.Parse(body)
	if err != nil {
		writeErrorResponse(rw, http.StatusInternalServerError, "Body could not be unmarshalled into a user struct", trace.Guarantee(err).Add("body", body))
		return
	}

	if err := user.EnsureFields(partialUser); err != nil {
		writeErrorResponse(rw, http.StatusInternalServerError, "User is missing required fields", trace.Guarantee(err).Add("user", partialUser))
		return
	}

	if partialUser.Id != "" || partialUser.CreatedAt != "" || partialUser.UpdatedAt != "" {
		writeErrorResponse(rw, http.StatusBadRequest, "Can only specify email and user_name when creating a user", ErrorInvalidBody.Add("user", partialUser))
		return
	}

	newUser := user.NewUser(partialUser.Email, partialUser.UserName)
	logrus.WithField("user", newUser).Debug("creating user")
	if _, err := user.CreateUser(config.Values.UsertableName, newUser); err != nil {
		writeErrorResponse(rw, http.StatusInternalServerError, "Could not create user", trace.Guarantee(err).Add("user", newUser))
		return
	}

	writeSuccessResponse(rw, newUser)
}

func getUser(rw http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("id") {
		writeErrorResponse(rw, http.StatusBadRequest, "Query must include an id", ErrorInvalidQuery.Add("query", r.URL.RawQuery))
		return
	}

	searchUser := &user.User{Id: r.URL.Query().Get("id")}
	foundUser, err := user.GetUser(config.Values.UsertableName, searchUser)
	if err != nil {
		writeErrorResponse(rw, http.StatusInternalServerError, "Error searching for user", trace.Guarantee(err).Add("user", searchUser))
		return
	}

	if foundUser == nil {
		writeErrorResponse(rw, http.StatusNotFound, "A matching user could not be found", ErrorNotFound.Add("user", searchUser))
		return
	}

	writeSuccessResponse(rw, foundUser)
}

func updateUser(rw http.ResponseWriter, r *http.Request) {
	if r.ContentLength < 1 {
		writeErrorResponse(rw, http.StatusBadRequest, "Body cannot be empty", ErrorInvalidBody)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(rw, http.StatusInternalServerError, "Body could not be read", trace.Guarantee(err).Add("body", body))
		return
	}

	partialUser, err := user.Parse(body)
	if err != nil {
		writeErrorResponse(rw, http.StatusInternalServerError, "Body could not be unmarshalled into a user struct", trace.Guarantee(err).Add("body", body))
		return
	}

	if partialUser.Id == "" {
		writeErrorResponse(rw, http.StatusBadRequest, "Must specify a user id", ErrorInvalidBody.Add("query", r.Body))
		return
	}

	updatedUser, err := user.UpdateUser(config.Values.UsertableName, partialUser)
	if err != nil {
		writeErrorResponse(rw, http.StatusInternalServerError, "Error updating user", trace.Guarantee(err).Add("user", partialUser))
		return
	}

	writeSuccessResponse(rw, updatedUser)
}
