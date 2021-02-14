package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

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
	answer := GetUserInput("> ")

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
}

func (c *Client) AddCredentials() {
	c.username = GetUserInput("Username> ")
	c.password = GetUserInput("Password> ")
}

func GetUserInput(msg string) string {
	buf := bufio.NewReader(os.Stdin)

	fmt.Print(msg)
	b, _ := buf.ReadBytes('\n')
	input := string(b)
	input = strings.TrimSuffix(input, "\n")
	return input
}

func (c *Client) Welcome() {
	// todo help
	for {
		action := GetUserInput(c.username + "> ")

		switch action {
		case "add_friend":
			friend := GetUserInput("Friend`s name> ")
			c.AddFriend(friend)
		case "create_group":
			break
		case "show_users":
			c.ShowUsers()
		case "owe":
			break
		case "lend":
			break
		case "payed":
			break // amount to reason ..
		case "payed_group":
			break // amount to reason ..
		case "exit":
			return
		default:
			fmt.Println("Invalid option. Try again!")
		}
	}
}

//
