package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Client struct {
	http.Client
	username string
	password string
}

// # Login or CreateAccount
// Login
func (c *Client) Login() error {
	req, _ := http.NewRequest("POST", UrlLogin, nil)
	req.SetBasicAuth(c.username, c.password)
	res, err := c.Do(req)
	if err != nil {
		return errors.New("ops, we couldn't process this")
	}

	if res.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(b))
	}
	return nil
}

// CreateAccount
func (c *Client) CreateAccount() error {
	req, _ := http.NewRequest("POST", UrlCreateAccount, nil)
	req.SetBasicAuth(c.username, c.password)
	res, err := c.Do(req)
	if err != nil {
		return errors.New("oops, we couldn't process that")
	}

	if res.StatusCode != http.StatusCreated {
		b, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(b))
	}
	return nil
}

// # Users
func (c *Client) ShowUsers() []string {
	req, _ := http.NewRequest("GET", UrlShowUsers, nil)
	req.SetBasicAuth(c.username, c.password)
	res, _ := c.Do(req)

	var u []string
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &u)
	return u
}

// # Friends
func (c *Client) AddFriend(friend string) error {
	reqBody, err := json.Marshal(map[string]string{
		"friend": friend,
	})
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", UrlAddFriend, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-type", "application/json")
	req.SetBasicAuth(c.username, c.password)
	res, err := c.Do(req)
	if err != nil {
		return errors.New("oops, we couldn't process that")
	}

	if res.StatusCode != http.StatusCreated {
		b, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(b))
	}
	return nil
}

func (c *Client) AddDebtToFriend(friend string, amount int, reason string, creditor bool) error {
	reqBody, err := json.Marshal(map[string]interface{}{
		"friend":   friend,
		"amount":   amount,
		"reason":   reason,
		"creditor": creditor,
	})
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", UrlAddDebt, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-type", "application/json")
	req.SetBasicAuth(c.username, c.password)
	res, err := c.Do(req)
	if err != nil {
		return errors.New("oops, we could not process that")
	}

	if res.StatusCode != http.StatusCreated {
		b, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(b))
	}
	return nil
}

func (c *Client) ShowFriends() []string {
	req, _ := http.NewRequest("GET", UrlShowFriends, nil)
	req.SetBasicAuth(c.username, c.password)
	res, _ := c.Do(req)

	var u []string
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &u)

	return u
}

func (c *Client) ShowDebts() []DebtC {
	req, _ := http.NewRequest("GET", UrlShowDebts, nil)
	req.SetBasicAuth(c.username, c.password)
	res, _ := c.Do(req)

	var owed []DebtC
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &owed)
	return owed
}

func (c *Client) ShowLoans() []DebtC {
	req, _ := http.NewRequest("GET", UrlShowLoans, nil)
	req.SetBasicAuth(c.username, c.password)
	res, _ := c.Do(req)

	var lent []DebtC
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &lent)
	return lent
}

// # Group
func (c *Client) CreateGroup(name string, participants []string) error {
	reqBody, err := json.Marshal(map[string]interface{}{
		"name":         name,
		"participants": participants,
	})

	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", UrlCreateGroup, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-type", "application/json")
	req.SetBasicAuth(c.username, c.password)
	res, err := c.Do(req)
	if err != nil {
		return errors.New("oops, we couldn't process that")
	}

	if res.StatusCode != http.StatusCreated {
		b, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(b))
	}
	return nil
}

func (c *Client) AddDebtToGroup(group string, amount int, reason string) error {
	reqBody, err := json.Marshal(map[string]interface{}{
		"group":  group,
		"amount": amount,
		"reason": reason,
	})
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", UrlAddDebtToGroup, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-type", "application/json")
	req.SetBasicAuth(c.username, c.password)
	res, err := c.Do(req)
	if err != nil {
		return errors.New("oops, we couldn't process that")
	}

	if res.StatusCode != http.StatusCreated {
		b, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(b))
	}
	return nil
}

func (c *Client) ReturnDebt(friend string, amount int, groupName string) error {
	reqBody, err := json.Marshal(map[string]interface{}{
		"friend":    friend,
		"amount":    amount,
		"groupName": groupName,
	})
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", UrlReturnDebt, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-type", "application/json")
	req.SetBasicAuth(c.username, c.password)
	res, err := c.Do(req)
	if err != nil {
		return errors.New("oops, we couldn't process that")
	}

	if res.StatusCode != http.StatusCreated {
		b, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(b))
	}
	return nil
}

func (c *Client) ShowGroups() []string {
	req, _ := http.NewRequest("GET", UrlShowGroups, nil)
	req.SetBasicAuth(c.username, c.password)
	res, _ := c.Do(req)

	var g []string
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &g)

	return g
}

func (c *Client) ShowDebtsToGroups() map[string][]DebtC {
	req, _ := http.NewRequest("GET", UrlShowDebtsToGroups, nil)
	req.SetBasicAuth(c.username, c.password)
	res, _ := c.Do(req)

	var owed map[string][]DebtC
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &owed)
	return owed
}

func (c *Client) ShowLoansToGroups() []DebtC {
	req, _ := http.NewRequest("GET", UrlShowLoansToGroups, nil)
	req.SetBasicAuth(c.username, c.password)
	res, _ := c.Do(req)

	var lent []DebtC
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &lent)
	return lent
}

type DebtC struct {
	To     string
	Amount int
	Reason string
}

func main() {
	var c Client
	c.Index()
}
