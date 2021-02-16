package main

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
)

const (
	password = "123456"
)

func TestApp_CreateAccount(t *testing.T) {
	t.Run("when successful", func(t *testing.T) {
		request, _ := http.NewRequest("POST", UrlCreateAccount, nil)
		request.SetBasicAuth("new", password)
		recorder := httptest.NewRecorder()
		app := InitApp()

		app.CreateAccount(recorder, request)

		response := recorder.Result()
		assert.Equal(t, http.StatusCreated, response.StatusCode)
	})
	t.Run("when creating user with existing username", func(t *testing.T) {
		request, _ := http.NewRequest("POST", UrlCreateAccount, nil)
		request.SetBasicAuth(peter, password)
		recorder := httptest.NewRecorder()
		app := InitApp()

		app.CreateAccount(recorder, request)

		response := recorder.Result()
		err, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, "this username is already taken\n", string(err))
	})
}

func TestApp_Login(t *testing.T) {
	t.Run("when successful", func(t *testing.T) {
		request, _ := http.NewRequest("POST", UrlCreateAccount, nil)
		request.SetBasicAuth(peter, peterPass)
		recorder := httptest.NewRecorder()
		app := InitApp()

		app.Login(recorder, request)

		response := recorder.Result()
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
	t.Run("when wrong credentials", func(t *testing.T) {
		request, _ := http.NewRequest("POST", UrlCreateAccount, nil)
		request.SetBasicAuth(peter, password)
		recorder := httptest.NewRecorder()
		app := InitApp()

		app.Login(recorder, request)

		response := recorder.Result()
		err, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
		assert.Equal(t, "invalid username or password\n", string(err))
	})
}

func TestApp_ShowUsers(t *testing.T) {
	t.Run("when having users", func(t *testing.T) {
		request, _ := http.NewRequest("GET", UrlShowUsers, nil)
		request.Header.Set("Username", peter)
		recorder := httptest.NewRecorder()
		app := InitApp()

		app.ShowUsers(recorder, request)

		response := recorder.Result()
		body, _ := ioutil.ReadAll(response.Body)

		actual := []string{}
		_ = json.Unmarshal(body, &actual)
		sort.Strings(actual)

		expected := []string{peter, george, lily, maria}
		sort.Strings(expected)

		assert.Equal(t, expected, actual)
		assert.Equal(t, 200, response.StatusCode)
	})
}

func TestApp_AddFriend(t *testing.T) {
	t.Run("when successful", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"friend": lily})
		request, _ := http.NewRequest("POST", UrlAddFriend, bytes.NewBuffer(body))
		request.Header.Set("Username", george)
		request.Header.Set("Content-type", "application/json")
		recorder := httptest.NewRecorder()
		app := InitApp()

		app.AddFriend(recorder, request)

		response := recorder.Result()
		assert.Equal(t, http.StatusCreated, response.StatusCode)
	})
	t.Run("when already friends", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{"friend": peter})
		request, _ := http.NewRequest("POST", UrlAddFriend, bytes.NewBuffer(body))
		request.Header.Set("Username", george)
		request.Header.Set("Content-type", "application/json")
		recorder := httptest.NewRecorder()
		app := InitApp()

		app.AddFriend(recorder, request)

		response := recorder.Result()
		err, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, "you are already friends\n", string(err))
	})
}

func TestApp_ShowFriends(t *testing.T) {
	t.Run("when having friends", func(t *testing.T) {
		request, _ := http.NewRequest("GET", UrlShowFriends, nil)
		request.Header.Set("Username", peter)
		recorder := httptest.NewRecorder()
		app := InitApp()

		app.ShowFriends(recorder, request)

		response := recorder.Result()
		body, _ := ioutil.ReadAll(response.Body)

		actual := []string{}
		_ = json.Unmarshal(body, &actual)
		sort.Strings(actual)

		expected := []string{george, lily}
		sort.Strings(expected)

		assert.Equal(t, expected, actual)
		assert.Equal(t, 200, response.StatusCode)
	})
	t.Run("when having no friends", func(t *testing.T) {
		request, _ := http.NewRequest("GET", UrlShowFriends, nil)
		request.Header.Set("Username", maria)
		recorder := httptest.NewRecorder()
		app := InitApp()

		app.ShowFriends(recorder, request)

		response := recorder.Result()
		body, _ := ioutil.ReadAll(response.Body)

		actual := []string{}
		_ = json.Unmarshal(body, &actual)

		expected := []string{}

		assert.Equal(t, expected, actual)
		assert.Equal(t, 200, response.StatusCode)
	})
}

