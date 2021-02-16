package main

import (
	"Cost-Sharing/server"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Launch server..")
	app := server.InitApp()

	// # Create or Login
	app.Server.Router.HandleFunc(server.PathToCreateAccount, app.CreateAccount).Methods("POST")
	app.Server.Router.HandleFunc(server.PathToLogin, app.Login).Methods("POST")

	// # Users
	app.Server.Router.HandleFunc(server.PathToShowUsers, server.Notify(&app, app.ShowUsers)).Methods("GET")

	// # Friends
	app.Server.Router.HandleFunc(server.PathToAddFriend, server.Notify(&app, app.AddFriend)).Methods("POST")
	app.Server.Router.HandleFunc(server.PathToShowFriends, server.Notify(&app, app.ShowFriends)).Methods("GET")

	// # Groups
	app.Server.Router.HandleFunc(server.PathToCreateGroup, server.Notify(&app, app.CreateGroup)).Methods("POST")
	app.Server.Router.HandleFunc(server.PathToShowGroups, server.Notify(&app, app.ShowGroups)).Methods("GET")

	// # Debt
	// ## To friends
	app.Server.Router.HandleFunc(server.PathToAddDebt, server.Notify(&app, app.AddDebtToFriend)).Methods("POST")

	app.Server.Router.HandleFunc(server.PathToShowDebts, server.Notify(&app, app.ShowDebts)).Methods("GET")
	app.Server.Router.HandleFunc(server.PathToShowLoans, server.Notify(&app, app.ShowLoans)).Methods("GET")

	// ## To group
	app.Server.Router.HandleFunc(server.PathToAddDebtToGroup, server.Notify(&app, app.AddDebtToGroup)).Methods("POST")
	app.Server.Router.HandleFunc(server.PathToReturnDebt, server.Notify(&app, app.ReturnDebt)).Methods("POST")

	app.Server.Router.HandleFunc(server.PathToShowDebtsToGroups, server.Notify(&app, app.ShowDebtsToGroups)).Methods("GET")
	app.Server.Router.HandleFunc(server.PathToShowLoansToGroups, server.Notify(&app, app.ShowLoansToGroups)).Methods("GET")

	log.Fatal(http.ListenAndServe(server.Port, app.Server.Router))
}
