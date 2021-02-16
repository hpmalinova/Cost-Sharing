package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func (c *Client) Index() {
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
	answer := getUserInput("> ")

	switch answer {
	case "login":
		c.AddCredentials()
		return c.Login()
	case "create_account":
		c.AddCredentials()
		return c.CreateAccount()
	case "exit":
		return nil
	default:
		return errors.New("invalid operation, try again")
	}
}

func (c *Client) AddCredentials() {
	c.username = getUserInput("Username> ")
	c.password = getUserInput("Password> ")
}

func (c *Client) Welcome() {
	// todo help
	for {
		action := getUserInput(c.username + "> ")

		switch action {
		case "show_users":
			c.showUsers()
		case "add_friend":
			c.addFriend()
		case "show_friends":
			c.showFriends()
		case "create_group":
			c.createGroup()
		case "show_groups":
			c.showGroups()
		case "split":
			c.split()
		case "split_group":
			c.splitGroup()
		case "pay_back":
			c.payBack()
		case "pay_back_group":
			c.payBackGroup()
		case "show_debts":
			c.showDebts()
		case "show_loans":
			c.showLoans()
		case "owe_group":
			break
		case "lend_group":
			break
		case "exit":
			return
		default:
			fmt.Println("Invalid option. Try again!")
		}
	}
}

func (c *Client) showUsers() {
	users := c.ShowUsers()
	printData(users, "Users: ", "There are no users!")
}

func (c *Client) addFriend() {
	friend := getUserInput("Friend`s name> ")
	err := c.AddFriend(friend)
	if err != nil {
		fmt.Println(err)
	} else {
		msg := "Congrats, you and " + friend + "are now friends!"
		fmt.Println(msg)
	}
}

func (c *Client) showFriends() {
	friends := c.ShowFriends()
	printData(friends, "Friends: ", "You have no friends!")
}

func (c *Client) createGroup() {
	groupName := getUserInput("Group`s name> ")

	participants := strings.Split(getUserInput("Participants (with `,`)> "), ",")
	for i, _ := range participants {
		participants[i] = strings.Trim(participants[i], " ")
	}

	err := c.CreateGroup(groupName, participants)
	if err != nil {
		fmt.Println(err)
	}
}

func (c *Client) showGroups() {
	groups := c.ShowGroups()
	printData(groups, "You participate in: ", "You don`t participate in any group!")
}

func (c *Client) split() {
	friend := getUserInput("Friend`s name> ")
	textAmount := getUserInput("Amount> ")
	amount, err := strconv.Atoi(textAmount)
	if err != nil || amount <= 0 {
		fmt.Println("Amount should be a number, bigger than 1!")
		return
	}
	reason := getUserInput("Reason for payment> ")
	splitAmount := int(math.Ceil(float64(amount) / 2))
	err = c.AddDebtToFriend(friend, splitAmount, reason, true)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (c *Client) splitGroup() {
	group := getUserInput("Group`s name> ")
	textAmount := getUserInput("Amount> ")
	amount, err := strconv.Atoi(textAmount)
	if err != nil || amount <= 0 {
		fmt.Println("Amount should be a number, bigger than 1!")
		return
	}
	reason := getUserInput("Reason for payment> ")
	err = c.AddDebtToGroup(group, amount, reason)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (c *Client) payBack() {
	// <username> payed me back <amount> lv
	friend := getUserInput("Friend`s name> ")
	textAmount := getUserInput("Amount> ")
	amount, err := strconv.Atoi(textAmount)
	if err != nil || amount <= 0 {
		fmt.Println("Amount should be a number, bigger than 1!")
		return
	}
	err = c.AddDebtToFriend(friend, amount, "", false)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (c *Client) payBackGroup() {
	// <username> payed me back <amount> lv for <groupName>
	friend := getUserInput("Friend`s name> ")
	textAmount := getUserInput("Amount> ")
	groupName := getUserInput("Group`s name> ")
	amount, err := strconv.Atoi(textAmount)
	if err != nil || amount <= 0 {
		fmt.Println("Amount should be a number, bigger than 1!")
		return
	}

	err = c.ReturnDebt(friend, amount, groupName)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (c *Client) showDebts() {
	o := c.ShowDebts()
	printDebt(o, "You owe money to: ", "You don`t owe any money!")
}

func (c *Client) showLoans() {
	l := c.ShowLoans()
	printDebt(l, "You`ve lent money to: ", "You haven`t lent any money!")
}

func getUserInput(msg string) string {
	buf := bufio.NewReader(os.Stdin)

	fmt.Print(msg)
	b, _ := buf.ReadBytes('\n')
	input := string(b)
	input = strings.TrimSuffix(input, "\n")
	return input
}

func printData(data []string, title, emptyMsg string) {
	if len(data) > 0 {
		fmt.Println(title)
		fmt.Println(strings.Join(data, ", "))
	} else {
		fmt.Println(emptyMsg)
	}
}

func printDebt(data []DebtC, title, empty string) {
	if len(data) > 0 {
		fmt.Println(title)
		for i, d := range data {
			_, _ = fmt.Fprintln(os.Stdout, i+1, ".", d.To, "-", d.Amount, "lv", "for", d.Reason)
		}
	} else {
		fmt.Println(empty)
	}
}
