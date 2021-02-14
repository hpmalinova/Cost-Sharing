package main

import (
	"Cost-Sharing/storage"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	//"golang.org/x/crypto/ssh/terminal"
	//_ "golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	//"syscall"
)

type Client struct {
	http.Client
	username string
	password string
}

func (c *Client) Authenticate() {
	fmt.Println("Welcome to Cost-Sharing")
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

func (c *Client) LoginOrCreate() error {
	fmt.Println("Do you want to create account or to login?")
	fmt.Println("Type `login`, `create_account` or `exit`")

	//for {
	answer := UserInput("> ")

	switch answer {
	case "login":
		c.AddCredentials()
		return c.Login()
	case "create_account":
		c.AddCredentials()
		return c.CreateUser()
	case "exit":
		return nil
	default:
		return errors.New("invalid operation, try again")
	}
	//if answer == "login" {
	//	c.AddCredentials()
	//	return c.Login()
	//} else if answer == "create_account" {
	//	c.AddCredentials()
	//	return c.CreateUser()
	//} else {
	//	fmt.Println("Type `login` or `create_account`")
	//}
	//}
}

func (c *Client) AddCredentials() {
	c.username = UserInput("Username> ")
	c.password = UserInput("Password> ")
}

func UserInput(msg string) string {
	buf := bufio.NewReader(os.Stdin)

	fmt.Print(msg)
	b, _ := buf.ReadBytes('\n')
	input := string(b)
	input = strings.TrimSuffix(input, "\n")
	return input
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

func (c *Client) Welcome() {
	// todo help
	for {
		action := UserInput(c.username + "> ")

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
		case "exit":
			return
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

func main() {
	var c Client
	c.Authenticate()
	//c.CreateUser("Pesho", "123456")
	//c.CreateUser("Silvia", "qwerty")
	//c.CreateUser("Silvia", "qwerty")
	//c.ShowUsers()
}
