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

func main() {
	var c Client
	c.username = "p"
	c.password = "1"
	c.CreateUser()

	c.username = "o"
	c.password = "1"
	c.CreateUser()

	c.Authenticate()
}
