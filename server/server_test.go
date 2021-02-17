package main

import (
	"Cost-Sharing/storage"
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
		app := InitApp(true)

		app.CreateAccount(recorder, request)

		response := recorder.Result()
		assert.Equal(t, http.StatusCreated, response.StatusCode)
	})
	t.Run("when creating user with existing username", func(t *testing.T) {
		request, _ := http.NewRequest("POST", UrlCreateAccount, nil)
		request.SetBasicAuth(peter, password)
		recorder := httptest.NewRecorder()
		app := InitApp(true)

		app.CreateAccount(recorder, request)

		response := recorder.Result()
		err, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, "this username is already taken\n", string(err))
	})
}

func TestApp_Login(t *testing.T) {
	t.Run("when successful", func(t *testing.T) {
		request, _ := http.NewRequest("POST", UrlLogin, nil)
		request.SetBasicAuth(peter, peterPass)
		recorder := httptest.NewRecorder()
		app := InitApp(true)

		app.Login(recorder, request)

		response := recorder.Result()
		assert.Equal(t, http.StatusOK, response.StatusCode)
	})
	t.Run("when wrong credentials", func(t *testing.T) {
		request, _ := http.NewRequest("POST", UrlLogin, nil)
		request.SetBasicAuth(peter, password)
		recorder := httptest.NewRecorder()
		app := InitApp(true)

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
		app := InitApp(true)

		app.ShowUsers(recorder, request)

		response := recorder.Result()
		body, _ := ioutil.ReadAll(response.Body)

		actual := []string{}
		_ = json.Unmarshal(body, &actual)
		sort.Strings(actual)

		expected := []string{peter, george, lily, maria, juan}
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
		app := InitApp(true)

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
		app := InitApp(true)

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
		app := InitApp(true)

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
		app := InitApp(true)

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

func TestApp_AddDebtToFriend(t *testing.T) {
	t.Run("when successful", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"friend":   george,
			"amount":   20,
			"reason":   "food",
			"creditor": true,
		})

		request, _ := http.NewRequest("POST", UrlAddDebt, bytes.NewBuffer(body))
		request.Header.Set("Username", peter)
		request.Header.Set("Content-type", "application/json")
		recorder := httptest.NewRecorder()
		app := InitApp(true)

		app.AddDebtToFriend(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusCreated, response.StatusCode)
	})
	t.Run("when not friends", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"friend":   george,
			"amount":   20,
			"reason":   "food",
			"creditor": true,
		})

		request, _ := http.NewRequest("POST", UrlAddDebt, bytes.NewBuffer(body))
		request.Header.Set("Username", lily)
		request.Header.Set("Content-type", "application/json")
		recorder := httptest.NewRecorder()
		app := InitApp(true)

		app.AddDebtToFriend(recorder, request)

		response := recorder.Result()
		err, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		assert.Equal(t, george+" is not your friend!\n", string(err))
	})
}

func TestApp_ShowDebts(t *testing.T) {
	t.Run("when having debts", func(t *testing.T) {
		request, _ := http.NewRequest("GET", UrlShowDebts, nil)
		request.Header.Set("Username", peter)
		recorder := httptest.NewRecorder()
		app := InitApp(true)

		app.ShowDebts(recorder, request)

		response := recorder.Result()
		body, _ := ioutil.ReadAll(response.Body)

		actual := []storage.DebtC{}
		_ = json.Unmarshal(body, &actual)

		expected := []storage.DebtC{{george, amount, food}}

		assert.True(t, storage.ContainsAll(expected, actual))
		assert.Equal(t, 200, response.StatusCode)
	})
}

func TestApp_ShowLoans(t *testing.T) {
	t.Run("when having loans", func(t *testing.T) {
		request, _ := http.NewRequest("GET", UrlShowLoans, nil)
		request.Header.Set("Username", peter)
		recorder := httptest.NewRecorder()
		app := InitApp(true)

		app.ShowLoans(recorder, request)

		response := recorder.Result()
		body, _ := ioutil.ReadAll(response.Body)

		actual := []storage.DebtC{}
		_ = json.Unmarshal(body, &actual)

		expected := []storage.DebtC{{lily, amount, food}}

		assert.True(t, storage.ContainsAll(expected, actual))
		assert.Equal(t, 200, response.StatusCode)
	})
}

func TestApp_CreateGroup(t *testing.T) {
	t.Run("when successful", func(t *testing.T) {
		participants := []string{lily, george}
		body, _ := json.Marshal(map[string]interface{}{
			"name":         "Party",
			"participants": participants,
		})
		request, _ := http.NewRequest("POST", UrlCreateGroup, bytes.NewBuffer(body))
		request.Header.Set("Username", george)
		request.Header.Set("Content-type", "application/json")
		recorder := httptest.NewRecorder()
		app := InitApp(true)

		app.CreateGroup(recorder, request)

		response := recorder.Result()
		assert.Equal(t, http.StatusCreated, response.StatusCode)
	})
}

func TestApp_ShowGroups(t *testing.T) {
	t.Run("when participating in groups", func(t *testing.T) {
		request, _ := http.NewRequest("GET", UrlShowGroups, nil)
		request.Header.Set("Username", peter)
		recorder := httptest.NewRecorder()
		app := InitApp(true)

		app.ShowGroups(recorder, request)

		response := recorder.Result()
		body, _ := ioutil.ReadAll(response.Body)

		actual := []string{}
		_ = json.Unmarshal(body, &actual)
		sort.Strings(actual)

		expected := []string{presents, japan, test}
		sort.Strings(expected)

		assert.Equal(t, expected, actual)
		assert.Equal(t, 200, response.StatusCode)
	})
	t.Run("when having no groups", func(t *testing.T) {
		request, _ := http.NewRequest("GET", UrlShowFriends, nil)
		request.Header.Set("Username", juan)
		recorder := httptest.NewRecorder()
		app := InitApp(true)

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

func TestApp_AddDebtToGroup(t *testing.T) {
	t.Run("when successful", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"group":  test,
			"amount": 30,
			"reason": food,
		})

		request, _ := http.NewRequest("POST", UrlAddDebtToGroup, bytes.NewBuffer(body))
		request.Header.Set("Username", peter)
		request.Header.Set("Content-type", "application/json")
		recorder := httptest.NewRecorder()
		app := InitApp(true)

		app.AddDebtToGroup(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusCreated, response.StatusCode)
	})
	t.Run("when client not a member", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"group":  test,
			"amount": 30,
			"reason": food,
		})

		request, _ := http.NewRequest("POST", UrlAddDebtToGroup, bytes.NewBuffer(body))
		request.Header.Set("Username", george)
		request.Header.Set("Content-type", "application/json")
		recorder := httptest.NewRecorder()
		app := InitApp(true)

		app.AddDebtToGroup(recorder, request)

		response := recorder.Result()
		err, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		msg := "you don`t participate in group called " + test + "\n"
		assert.Equal(t, msg, string(err))
	})

}

func TestApp_ReturnDebt(t *testing.T) {
	t.Run("when successful", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"friend":    peter,
			"amount":    amount,
			"groupName": presents,
		})

		request, _ := http.NewRequest("POST", UrlReturnDebt, bytes.NewBuffer(body))
		request.Header.Set("Username", lily)
		request.Header.Set("Content-type", "application/json")
		recorder := httptest.NewRecorder()
		app := InitApp(true)

		app.ReturnDebt(recorder, request)

		response := recorder.Result()

		assert.Equal(t, http.StatusCreated, response.StatusCode)
	})
	t.Run("when client not a member", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"friend":    peter,
			"amount":    amount,
			"groupName": presents,
		})

		request, _ := http.NewRequest("POST", UrlReturnDebt, bytes.NewBuffer(body))
		request.Header.Set("Username", george)
		request.Header.Set("Content-type", "application/json")
		recorder := httptest.NewRecorder()
		app := InitApp(true)

		app.ReturnDebt(recorder, request)

		response := recorder.Result()
		err, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		msg := "you don`t participate in group called " + presents + "\n"
		assert.Equal(t, msg, string(err))
	})
	t.Run("when friend not a member", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"friend":    george,
			"amount":    amount,
			"groupName": presents,
		})

		request, _ := http.NewRequest("POST", UrlReturnDebt, bytes.NewBuffer(body))
		request.Header.Set("Username", lily)
		request.Header.Set("Content-type", "application/json")
		recorder := httptest.NewRecorder()
		app := InitApp(true)

		app.ReturnDebt(recorder, request)

		response := recorder.Result()
		err, _ := ioutil.ReadAll(response.Body)

		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
		msg := george + " is not part of " + presents + "\n"
		assert.Equal(t, msg, string(err))
	})
}

func TestApp_ShowDebtsToGroups(t *testing.T) {
	t.Run("when having debts", func(t *testing.T) {
		request, _ := http.NewRequest("GET", UrlShowDebtsToGroups, nil)
		request.Header.Set("Username", lily)
		recorder := httptest.NewRecorder()
		app := InitApp(true)

		app.ShowDebtsToGroups(recorder, request)

		response := recorder.Result()
		body, _ := ioutil.ReadAll(response.Body)

		actual := map[string][]storage.DebtC{}
		_ = json.Unmarshal(body, &actual)

		expected := map[string][]storage.DebtC{presents: {{peter, amount, food}}}

		assert.True(t, storage.Equal(expected, actual))
		assert.Equal(t, 200, response.StatusCode)
	})
}

func TestApp_ShowLoansToGroups(t *testing.T) {
	t.Run("when having loans", func(t *testing.T) {
		request, _ := http.NewRequest("GET", UrlShowLoansToGroups, nil)
		request.Header.Set("Username", peter)
		recorder := httptest.NewRecorder()
		app := InitApp(true)

		app.ShowLoansToGroups(recorder, request)

		response := recorder.Result()
		body, _ := ioutil.ReadAll(response.Body)

		actual := map[string][]storage.DebtC{}
		_ = json.Unmarshal(body, &actual)

		expected := map[string][]storage.DebtC{
			presents: {{lily, amount, food}},
			japan:    {{george, amount, food}, {maria, amount, food}},
			test:     {},
		}

		assert.True(t, storage.Equal(expected, actual))
		assert.Equal(t, 200, response.StatusCode)
	})
}
