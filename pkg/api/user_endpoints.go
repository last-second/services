package api

import (
	"io"
	"net/http"

	"github.com/last-second/services/pkg/config"
	"github.com/last-second/services/pkg/user"
	"github.com/sirupsen/logrus"
)

func createUser(rw http.ResponseWriter, r *http.Request) {
	if r.ContentLength < 1 {
		writeErrorResponse(
			rw, http.StatusBadRequest,
			"Body cannot be empty",
			ErrInvalidBody, nil,
		)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(
			rw, http.StatusInternalServerError,
			"Body could not be read",
			ErrParseBody.Add("body", body), err,
		)
		return
	}

	partialUser, err := user.Parse(body)
	if err != nil {
		writeErrorResponse(
			rw, http.StatusInternalServerError,
			"Body could not be unmarshalled into a user struct",
			ErrInvalidBody.Add("body", body), err,
		)
		return
	}

	if partialUser.Id != "" || partialUser.CreatedAt != "" || partialUser.UpdatedAt != "" {
		writeErrorResponse(
			rw, http.StatusBadRequest,
			"Can only specify email and user_name when creating a user",
			ErrInvalidBody.Add("user", partialUser), nil,
		)
		return
	}

	newUser := user.NewUser(partialUser.Email, partialUser.UserName)
	logrus.WithField("user", newUser).Debug("creating user")
	if _, err := user.CreateUser(config.Values.UsertableName, newUser); err != nil {
		writeErrorResponse(
			rw, http.StatusInternalServerError,
			"Could not create user",
			ErrCreateUser.Add("user", newUser), err,
		)
		return
	}

	writeSuccessResponse(rw, newUser)
}

func getUser(rw http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("id") {
		writeErrorResponse(
			rw, http.StatusBadRequest,
			"Query must include an id",
			ErrInvalidQuery.Add("query", r.URL.RawQuery), nil,
		)
		return
	}

	searchUser := &user.User{
		Id:        r.URL.Query().Get("id"),
		Email:     r.URL.Query().Get("email"),
		UserName:  r.URL.Query().Get("user_name"),
		CreatedAt: r.URL.Query().Get("created_at"),
		UpdatedAt: r.URL.Query().Get("updated_at"),
	}

	foundUser, err := user.GetUser(config.Values.UsertableName, searchUser)
	if err != nil {
		writeErrorResponse(
			rw, http.StatusInternalServerError,
			"Error searching for user",
			ErrInvalidQuery.Add("user", searchUser), err,
		)
		return
	}

	if foundUser == nil {
		writeErrorResponse(
			rw, http.StatusNotFound,
			"A matching user could not be found",
			ErrNotFound.Add("user", searchUser), err,
		)
		return
	}

	writeSuccessResponse(rw, foundUser)
}
