package main

import (
	"Cost-Sharing/server"
	"net/http"
)

func main() {
	app := server.InitApp()

	app.Server.Router.HandleFunc("/costSharing", app.CreateUser).Methods("POST")
	app.Server.Router.HandleFunc("/costSharing/allUsers", app.ShowUsers).Methods("GET")

	_ = http.ListenAndServe(":8080", app.Server.Router)
}
