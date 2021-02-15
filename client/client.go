package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"

	//"golang.org/x/crypto/ssh/terminal"
	//_ "golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"net/http"
	//"syscall"
)

type Client struct {
	http.Client
	username string
	password string
}

func (c *Client) Login() error {
	req, _ := http.NewRequest("POST", "http://localhost:8080/costSharing/login", nil)
	req.SetBasicAuth(c.username, c.password)
	res, err := c.Do(req)
	if err != nil {
		return errors.New("ops, we couldn't process this") // todo
	}

	if res.StatusCode != 200 {
		b, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(b))
	}
	return nil
}

func (c *Client) CreateUser() error {
	req, _ := http.NewRequest("POST", "http://localhost:8080/costSharing/createAccount", nil)
	req.SetBasicAuth(c.username, c.password)
	res, err := c.Do(req)
	if err != nil {
		return errors.New("oops, we couldn't process that") // todo
	}

	if res.StatusCode != http.StatusCreated {
		b, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(b))
	}
	return nil
}

func (c *Client) ShowUsers() []string {
	req, _ := http.NewRequest("GET", "http://localhost:8080/costSharing/home/showUsers", nil)
	req.SetBasicAuth(c.username, c.password)
	res, _ := c.Do(req)
	var u []string
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &u)
	return u
}

func (c *Client) AddFriend(friend string) error {
	reqBody, err := json.Marshal(map[string]string{
		"friend": friend,
	})
	if err != nil {
		log.Println(err)
		return err // TODO
	}

	req, _ := http.NewRequest("POST", "http://localhost:8080/costSharing/home/addFriend", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-type", "application/json")
	req.SetBasicAuth(c.username, c.password)
	res, err := c.Do(req)
	if err != nil {
		log.Println(err)
		return errors.New("oops, we couldn't process that") // todo
	}

	if res.StatusCode != http.StatusCreated {
		b, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(b))
	}
	return nil
}

func (c *Client) ShowFriends() []string {
	req, _ := http.NewRequest("GET", "http://localhost:8080/costSharing/home/showFriends", nil)
	req.SetBasicAuth(c.username, c.password)
	res, _ := c.Do(req)

	var u []string
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &u)

	return u
	//fmt.Println("All friends: ", u)
}

func (c *Client) CreateGroup(name string, participants []string) error {
	reqBody, err := json.Marshal(map[string]interface{}{
		"name":         name,
		"participants": participants,
	})

	if err != nil {
		log.Println(err)
		return err // TODO
	}

	req, _ := http.NewRequest("POST", "http://localhost:8080/costSharing/home/createGroup", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-type", "application/json")
	req.SetBasicAuth(c.username, c.password)
	res, err := c.Do(req)
	if err != nil {
		log.Println(err)
		return errors.New("oops, we couldn't process that") // todo
	}

	if res.StatusCode != http.StatusCreated {
		b, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(b))
	}
	return nil
}

func (c *Client) ShowGroups() []string {
	req, _ := http.NewRequest("GET", "http://localhost:8080/costSharing/home/showGroups", nil)
	req.SetBasicAuth(c.username, c.password)
	res, _ := c.Do(req)

	var g []string
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &g)

	return g
}

func (c *Client) AddDebtToFriend(friend string, amount int, reason string, creditor bool) error {
	reqBody, err := json.Marshal(map[string]interface{}{
		"friend":   friend,
		"amount":   amount,
		"reason":   reason,
		"creditor": creditor,
	})
	if err != nil {
		log.Println(err)
		return err // TODO
	}

	req, _ := http.NewRequest("POST", "http://localhost:8080/costSharing/home/addDebt ", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-type", "application/json")
	req.SetBasicAuth(c.username, c.password)
	res, err := c.Do(req)
	if err != nil {
		log.Println(err)
		return errors.New("oops, we could not process that") // todo
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
		log.Println(err)
		return err // TODO
	}

	req, _ := http.NewRequest("POST", "http://localhost:8080/costSharing/home/addDebtGroup", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-type", "application/json")
	req.SetBasicAuth(c.username, c.password)
	res, err := c.Do(req)
	if err != nil {
		log.Println(err)
		return errors.New("oops, we couldn't process that") // todo
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
		log.Println(err)
		return err // TODO
	}

	req, _ := http.NewRequest("POST", "http://localhost:8080/costSharing/home/returnDebt", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-type", "application/json")
	req.SetBasicAuth(c.username, c.password)
	res, err := c.Do(req)
	if err != nil {
		log.Println(err)
		return errors.New("oops, we couldn't process that") // todo
	}

	if res.StatusCode != http.StatusCreated {
		b, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(b))
	}
	return nil
}

type DebtC struct {
	To     string
	Amount int
	Reason string
}

func (c *Client) ShowOwed() []DebtC {
	req, _ := http.NewRequest("GET", "http://localhost:8080/costSharing/home/owe", nil)
	req.SetBasicAuth(c.username, c.password)
	res, _ := c.Do(req)

	var owed []DebtC
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &owed)
	return owed
}

func (c *Client) ShowLent() []DebtC {
	req, _ := http.NewRequest("GET", "http://localhost:8080/costSharing/home/lend", nil)
	req.SetBasicAuth(c.username, c.password)
	res, _ := c.Do(req)

	var lent []DebtC
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &lent)
	return lent
}

func (c *Client) ShowOwedGroup() map[string][]DebtC {
	req, _ := http.NewRequest("GET", "http://localhost:8080/costSharing/home/oweGroup", nil)
	req.SetBasicAuth(c.username, c.password)
	res, _ := c.Do(req)

	var owed map[string][]DebtC
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &owed)
	return owed
}

func (c *Client) ShowLentGroup() []DebtC {
	req, _ := http.NewRequest("GET", "http://localhost:8080/costSharing/home/lendGroup", nil)
	req.SetBasicAuth(c.username, c.password)
	res, _ := c.Do(req)

	var lent []DebtC
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &lent)
	return lent
}

func main() {
	var c Client
	c.username = "p"
	c.password = "1"
	c.CreateUser()

	c.username = "o"
	c.password = "1"
	c.CreateUser()

	c.username = "r"
	c.password = "1"
	c.CreateUser()

	c.Index()
}
