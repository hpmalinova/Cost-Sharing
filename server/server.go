package server

import (
	"Cost-Sharing/storage"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"math"
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
	Users        storage.Users
	Friends      storage.Friends
	Money        storage.MoneyExchange
	Groups       storage.Groups
	Participants storage.Participants
	Participates storage.Participates
	Server       Server
}

func InitApp() App {
	return App{
		Users:   storage.Users{Users: map[storage.Username]storage.User{}},
		Friends: storage.Friends{Friends: map[storage.Username]map[storage.Username]struct{}{}},
		Money: storage.MoneyExchange{
			Owes:  map[storage.Username]storage.To{},
			Lends: map[storage.Username]storage.To{},
		},
		Groups:       storage.Groups{Groups: map[uuid.UUID]storage.Group{}},
		Participants: storage.Participants{Participants: map[uuid.UUID]map[storage.Username]struct{}{}},
		Participates: storage.Participates{Participates: map[storage.Username]map[uuid.UUID]struct{}{}},
		Server:       NewServer(),
	}
}

func Notify(a *App, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || !a.CheckPassword(username, password) {
			w.WriteHeader(401)
		} else {
			f(w, r)
		}
	}
}

func (a *App) CreateUser(res http.ResponseWriter, req *http.Request) {
	username, password, _ := req.BasicAuth()

	err := a.Users.Create(storage.Username(username), password)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusCreated)
}

func (a *App) Login (res http.ResponseWriter, req *http.Request) {
	username, password, _ := req.BasicAuth()
	err := a.Users.CheckCredentials(storage.Username(username), password)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
	}
}

func (a *App) CheckPassword(username, password string) bool {
	realPassword := a.Users.GetPassword(storage.Username(username))
	return CheckPasswordHash(password, realPassword)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}


func (a *App) ShowUsers(res http.ResponseWriter, req *http.Request) {
	marshal, _ := json.Marshal(a.Users.GetUsernames())
	_, _ = res.Write(marshal)
}

func (a *App) AddDebt(res http.ResponseWriter, req *http.Request) {
	// TODO where to put?
	// Given: creditor, groupName, amount, reason
	// TODO get from body
	var (
		creditor  storage.Username
		groupName string
		amount    int
		reason    string
	)

	// p.GetGroups(creditor) --> groupIDs //a.Participates.GetGroups()
	groupIDs := a.Participates.GetGroups(creditor)

	// check for every groupID
	//		if:	g.GetName(groupID) == groupName --> groupID
	var groupID uuid.UUID

	for _, id := range groupIDs {
		if a.Groups.GetName(id) == groupName {
			groupID = id
		}
	}

	// individualAmount := amount / p.GetNumberOfParticipants(groupID) // TODO math.Ceil(x float64)--> float64
	fAmount := math.Ceil(float64(amount) / float64(a.Participants.GetNumberOfParticipants(groupID)))
	debt := int(fAmount)

	// p.GetParticipants --> []Username
	participants := a.Participants.GetParticipants(groupID)

	// g.AddDebt(creditor, groupID, participants, individualAmount, reason)
	a.Groups.AddDebt(creditor, groupID, participants, debt, reason)
}
