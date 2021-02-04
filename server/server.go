package server

import (
	"Cost-Sharing/storage"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	Router *mux.Router
}

func NewServer() Server {
	return Server{
		Router: mux.NewRouter(),
	}
}

type App struct {
	Users   storage.Users
	Friends storage.Friends
	Money   storage.MoneyExchange
	Server  Server
}

func InitApp() App {
	return App{
		Users:   storage.Users{Users: map[storage.Username]storage.User{}},
		Friends: storage.Friends{Friends: map[storage.Username]map[storage.Username]struct{}{}},
		Money: storage.MoneyExchange{
			Owes:  map[storage.Username]storage.To{},
			Lends: map[storage.Username]storage.To{},
		},
		Server: NewServer(),
	}
}

func (a *App) CreateUser(res http.ResponseWriter, req *http.Request) {
	username, password, _ := req.BasicAuth()
	err := a.Users.Create(storage.Username(username), password)
	if err != nil {
		msg := fmt.Sprint(username + " is already taken")
		http.Error(res, msg, http.StatusForbidden)
	}
}

func (a *App) ShowUsers(res http.ResponseWriter, req *http.Request) {
	marshal, _ := json.Marshal(a.Users.GetUsernames())
	_, _ = res.Write(marshal)
}
