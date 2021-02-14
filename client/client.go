package main

import (
	"Cost-Sharing/storage"
	"encoding/json"
	"errors"
	"fmt"

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

	if res.StatusCode != 201 {
		b, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(b))
	}
	return nil
}

func (c *Client) ShowUsers() {
	req, _ := http.NewRequest("GET", "http://localhost:8080/costSharing/home/allUsers", nil)
	req.SetBasicAuth(c.username, c.password)
	res, _ := c.Do(req)
	var u []storage.Username
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &u)

	fmt.Println("All users: ", u)
}

func main() {
	var c Client
	c.username = "p"
	c.password = "1"
	c.CreateUser()

	c.Authenticate()
}
