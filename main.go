package main

import (
	"Cost-Sharing/server"
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Launch server..")
	app := server.InitApp()

	// # Create or Login
	app.Server.Router.HandleFunc("/costSharing/createAccount", app.CreateUser).Methods("POST")
	app.Server.Router.HandleFunc("/costSharing/login", app.Login).Methods("POST")

	//app.Server.Router.HandleFunc("/costSharing/home", app.Welcome).Methods("POST")

	// # Users
	app.Server.Router.HandleFunc("/costSharing/home/showUsers", server.Notify(&app, app.ShowUsers)).Methods("GET") // todo show users

	// # Friends
	app.Server.Router.HandleFunc("/costSharing/home/addFriend", server.Notify(&app, app.AddFriend)).Methods("POST")
	app.Server.Router.HandleFunc("/costSharing/home/showFriends", server.Notify(&app, app.ShowFriends)).Methods("GET")

	// # Groups
	app.Server.Router.HandleFunc("/costSharing/home/createGroup", server.Notify(&app, app.CreateGroup)).Methods("POST")
	//app.Server.Router.HandleFunc("/costSharing/home/showGroups", server.Notify(&app, app.ShowGroups)).Methods("GET")

	_ = http.ListenAndServe(":8080", app.Server.Router)
}
