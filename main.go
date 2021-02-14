package main

import (
	"Cost-Sharing/server"
	"net/http"
)

func main() {
	app := server.InitApp()

	app.Server.Router.HandleFunc("/costSharing/", app.Index).Methods("POST")

	//app.Server.Router.HandleFunc("/costSharing/createAccount", app.CreateUser).Methods("POST")
	//app.Server.Router.HandleFunc("/costSharing/login", app.Login).Methods("POST")

	//app.Server.Router.HandleFunc("/costSharing/home", app.Welcome).Methods("POST")

	app.Server.Router.HandleFunc("/costSharing/home/allUsers", server.Notify(&app, app.ShowUsers)).Methods("GET")

	_ = http.ListenAndServe(":8080", app.Server.Router)
}
