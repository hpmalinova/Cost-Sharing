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
		case "show_users":
			users := c.ShowUsers()
			printData(users, "Users: ", "There are no users!")
		case "add_friend":
			friend := GetUserInput("Friend`s name> ")
			err := c.AddFriend(friend)
			if err != nil {
				fmt.Println(err)
			} else {
				msg := "Congrats, you and " + friend + "are now friends!"
				fmt.Println(msg)
			}
		case "show_friends":
			friends := c.ShowFriends()
			printData(friends, "Friends: ", "You have no friends!")
		case "create_group":
			groupName := GetUserInput("Group`s name> ")

			participants := strings.Split(GetUserInput("Participants (with `,`)> "), ",")
			for i, _ := range participants {
				participants[i] = strings.Trim(participants[i], " ")
			}

			err := c.CreateGroup(groupName, participants)
			if err != nil {
				fmt.Println(err)
			} //} else {
			//	msg := "Congrats, you and " + friend + "are now friends!"
			//	fmt.Println(msg)
			//}
		case "show_groups":
			groups := c.ShowGroups()
			printData(groups, "You participate in: ", "You don`t participate in any group!")
		case "split":
			friend := GetUserInput("Friend`s name> ")
			textAmount := GetUserInput("Amount> ")
			amount, err := strconv.Atoi(textAmount)
			if err != nil || amount <= 0 {
				fmt.Println("Amount should be a number, bigger than 1!")
				continue
			}
			reason := GetUserInput("Reason for payment> ")
			splitAmount := int(math.Ceil(float64(amount) / 2))
			err = c.AddDebtToFriend(friend, splitAmount, reason, true)
			if err != nil {
				fmt.Println(err.Error())
			}
		case "split_group":
			group := GetUserInput("Group`s name> ")
			textAmount := GetUserInput("Amount> ")
			amount, err := strconv.Atoi(textAmount)
			if err != nil || amount <= 0 {
				fmt.Println("Amount should be a number, bigger than 1!")
				continue
			}
			reason := GetUserInput("Reason for payment> ")
			err = c.AddDebtToGroup(group, amount, reason)
			if err != nil {
				fmt.Println(err.Error())
			}
		case "payed_back":
			// <username> payed me back <amount> lv
			friend := GetUserInput("Friend`s name> ")
			textAmount := GetUserInput("Amount> ")
			amount, err := strconv.Atoi(textAmount)
			if err != nil || amount <= 0 {
				fmt.Println("Amount should be a number, bigger than 1!")
				continue
			}
			err = c.AddDebtToFriend(friend, amount, "", false)
			if err != nil {
				fmt.Println(err.Error())
			}
		case "pay_back_group":
			// <username> payed me back <amount> lv for <groupName>
			friend := GetUserInput("Friend`s name> ")
			textAmount := GetUserInput("Amount> ")
			groupName := GetUserInput("Group`s name> ")
			amount, err := strconv.Atoi(textAmount)
			if err != nil || amount <= 0 {
				fmt.Println("Amount should be a number, bigger than 1!")
				continue
			}
			// TODO
			err = c.ReturnDebt(friend, amount, groupName)
			if err != nil {
				fmt.Println(err.Error())
			}
		case "owe":
			o := c.ShowOwed()
			printDebt(o, "You owe money to: ", "You don`t owe any money!")
			// Silvia - 50 for food
		case "lend":
			l := c.ShowLent()
			printDebt(l, "You`ve lent money to: ", "You haven`t lent any money!")
		case "exit":
			return
		default:
			fmt.Println("Invalid option. Try again!")
		}
	}
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
