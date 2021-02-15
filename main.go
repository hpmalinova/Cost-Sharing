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

	// # Users
	app.Server.Router.HandleFunc("/costSharing/home/showUsers", server.Notify(&app, app.ShowUsers)).Methods("GET") // todo show users

	// # Friends
	app.Server.Router.HandleFunc("/costSharing/home/addFriend", server.Notify(&app, app.AddFriend)).Methods("POST")
	app.Server.Router.HandleFunc("/costSharing/home/showFriends", server.Notify(&app, app.ShowFriends)).Methods("GET")

	// # Groups
	app.Server.Router.HandleFunc("/costSharing/home/createGroup", server.Notify(&app, app.CreateGroup)).Methods("POST")
	app.Server.Router.HandleFunc("/costSharing/home/showGroups", server.Notify(&app, app.ShowGroups)).Methods("GET")

	// # Debt
	app.Server.Router.HandleFunc("/costSharing/home/split", server.Notify(&app, app.AddDebtToFriend)).Methods("POST")
	app.Server.Router.HandleFunc("/costSharing/home/split_group", server.Notify(&app, app.AddDebtToGroup)).Methods("POST")

	app.Server.Router.HandleFunc("/costSharing/home/owe", server.Notify(&app, app.ShowOwe)).Methods("GET")
	app.Server.Router.HandleFunc("/costSharing/home/lend", server.Notify(&app, app.AddDebtToGroup)).Methods("GET")

	_ = http.ListenAndServe(":8080", app.Server.Router)
}
