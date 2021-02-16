package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Launch server..")
	app := InitApp()

	// # Create or Login
	app.Server.Router.HandleFunc(PathToCreateAccount, app.CreateAccount).Methods("POST")
	app.Server.Router.HandleFunc(PathToLogin, app.Login).Methods("POST")

	// # Users
	app.Server.Router.HandleFunc(PathToShowUsers, Notify(&app, app.ShowUsers)).Methods("GET")

	// # Friends
	app.Server.Router.HandleFunc(PathToAddFriend, Notify(&app, app.AddFriend)).Methods("POST")
	app.Server.Router.HandleFunc(PathToShowFriends, Notify(&app, app.ShowFriends)).Methods("GET")

	app.Server.Router.HandleFunc(PathToAddDebt, Notify(&app, app.AddDebtToFriend)).Methods("POST")

	app.Server.Router.HandleFunc(PathToShowDebts, Notify(&app, app.ShowDebts)).Methods("GET")
	app.Server.Router.HandleFunc(PathToShowLoans, Notify(&app, app.ShowLoans)).Methods("GET")

	// # Groups
	app.Server.Router.HandleFunc(PathToCreateGroup, Notify(&app, app.CreateGroup)).Methods("POST")
	app.Server.Router.HandleFunc(PathToShowGroups, Notify(&app, app.ShowGroups)).Methods("GET")

	app.Server.Router.HandleFunc(PathToAddDebtToGroup, Notify(&app, app.AddDebtToGroup)).Methods("POST")
	app.Server.Router.HandleFunc(PathToReturnDebt, Notify(&app, app.ReturnDebt)).Methods("POST")

	app.Server.Router.HandleFunc(PathToShowDebtsToGroups, Notify(&app, app.ShowDebtsToGroups)).Methods("GET")
	app.Server.Router.HandleFunc(PathToShowLoansToGroups, Notify(&app, app.ShowLoansToGroups)).Methods("GET")

	log.Fatal(http.ListenAndServe(Port, app.Server.Router))
}
