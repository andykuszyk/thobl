package handlers

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"strings"
	"encoding/json"

	"github.com/stretchr/testify/assert"
	"github.com/andykuszyk/thobl/repos"
	"github.com/andykuszyk/thobl/models"
)

func TestPostUsersHandler_ShouldReturn200_WhenPostedUser(t *testing.T) {
	handler := BuildHandlers(&repos.UsersRepoMock{})
	req, err := http.NewRequest(http.MethodPost, "/users", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()

	handler.postUsersHandler(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPostUsersAuthenticateHandler_ShouldReturn404_WhenUserNotCreated(t *testing.T) {
	handler := BuildHandlers(&repos.UsersRepoMock{
		GetByUsernameFunc: func (s string) (*models.User, error) {
			return nil, nil
		},
	})
	req, err := http.NewRequest(http.MethodPost, "/users/authenticate", strings.NewReader(`{"username":"foo","password":"bar"}`))
	assert.NoError(t, err)
	w:= httptest.NewRecorder()

	handler.postUsersAuthenticateHandler(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPostUsersAuthenticateHandler_ShouldReturn201_WhenUserExists(t *testing.T) {
	handler := BuildHandlers(&repos.UsersRepoMock{
		GetByUsernameFunc: func (s string) (*models.User, error) {
			return &models.User{Username: "foo", Password: "bar"}, nil
		},
	})
	req, err := http.NewRequest(http.MethodPost, "/users/authenticate", strings.NewReader(`{"username":"foo","password":"bar"}`))
	assert.NoError(t, err)
	w:= httptest.NewRecorder()

	handler.postUsersAuthenticateHandler(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPostUsersAuthenticateHandler_ShouldReturnToken_WhenUserExists(t *testing.T) {
	handler := BuildHandlers(&repos.UsersRepoMock{
		GetByUsernameFunc: func (s string) (*models.User, error) {
			return &models.User{Username: "foo", Password: "bar"}, nil
		},
	})
	req, err := http.NewRequest(http.MethodPost, "/users/authenticate", strings.NewReader(`{"username":"foo","password":"bar"}`))
	assert.NoError(t, err)
	w:= httptest.NewRecorder()

	handler.postUsersAuthenticateHandler(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var authenticateResponse models.UserToken
	assert.NoError(t, json.NewDecoder(w.Body).Decode(&authenticateResponse))
	assert.NotEmpty(t, authenticateResponse.Token)
}

func TestPostUsersAuthenticateHandler_ShouldReturn401_WhenPasswordDoesntMatch(t *testing.T) {
	handler := BuildHandlers(&repos.UsersRepoMock{
		GetByUsernameFunc: func(s string) (*models.User, error) {
			return &models.User{
				Username: "foo",
				Password: "bar",
			}, nil
		},
	})
	req, err := http.NewRequest(http.MethodPost, "/users/authenticate", strings.NewReader(`{"username":"foo","password":"spam"}`))
	assert.NoError(t, err)
	w:= httptest.NewRecorder()

	handler.postUsersAuthenticateHandler(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
