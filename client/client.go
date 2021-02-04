package main

import (
	"Cost-Sharing/storage"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Client struct {
	http.Client
}

func (c *Client) CreateUser(username, password string) {
	req, _ := http.NewRequest("POST", "http://localhost:8080/costSharing", nil)
	req.SetBasicAuth(username, password)
	res, err := c.Do(req)
	if err != nil {
		log.Fatal(err) // TODO fix
	}
	if res.StatusCode != 200 {
		b, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(b))
	} else {
		fmt.Println("User created successfully") // TODO redirect
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
	c.CreateUser("Pesho", "123456")
	c.CreateUser("Silvia", "qwerty")
	c.CreateUser("Silvia", "qwerty")
	c.ShowUsers()
}
