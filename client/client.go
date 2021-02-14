package main

import (
	"Cost-Sharing/storage"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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
		return errors.New("ops, we couldn't process this") // todo
	}

	if res.StatusCode != 201 {
		b, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(b))
	}
	return nil
}

func (c *Client) Authenticate() {
	for {
		err := c.LoginOrCreate()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			break
		}
	}
	c.Welcome()
}

func (c *Client) Welcome() {
	// todo help
	for {
		buf := bufio.NewReader(os.Stdin)
		prompt := c.username + "> "
		fmt.Println(prompt)
		b, _ := buf.ReadBytes('\n')
		action := string(b)

		switch action {
		case "add_friend":
			break
		case "create_group":
			break
		case "show users":
			break
		case "owe":
			break
		case "lend":
			break
		case "payed":
			break // amount to reason ..
		case "payed group":
			break // amount to reason ..
		case "logout": // TODO ?
			c.username = ""
			c.password = ""
			c.Authenticate()
			break //fallthrough
		case "exit":
			break
		default:
			fmt.Println("Invalid option. Try again!")
		}
	}
}

func (c *Client) ShowUsers() {
	req, _ := http.NewRequest("GET", "http://localhost:8080/costSharing/allUsers", nil)
	res, _ := c.Do(req)

	var u []storage.Username
	body, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(body, &u)

	fmt.Println("All users: ", u)
}

func (c *Client) LoginOrCreate() error {
	buf := bufio.NewReader(os.Stdin)
	fmt.Println("Login or create account (l/c)")
	b, _ := buf.ReadBytes('\n')
	answer := string(b)

	switch answer {
	case "login":
		c.AddCredentials()
	case "create_account":
		c.AddCredentials()
		return c.CreateUser()
	default:
		return errors.New("try again")
	}
	return nil
}

func (c *Client) AddCredentials() { //todo add err?
	buf := bufio.NewReader(os.Stdin)
	fmt.Println("Username:")
	b, _ := buf.ReadBytes('\n')
	username := string(b)
	c.username = username

	buf = bufio.NewReader(os.Stdin)
	fmt.Println("Username:")
	b, _ = buf.ReadBytes('\n')
	password := string(b)
	c.password = password
}

func main() {
	//var c Client
	//c.CreateUser("Pesho", "123456")
	//c.CreateUser("Silvia", "qwerty")
	//c.CreateUser("Silvia", "qwerty")
	//c.ShowUsers()
}
