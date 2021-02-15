package server

import (
	"Cost-Sharing/storage"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
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
		Users:   storage.Users{Users: map[string]storage.User{}},
		Friends: storage.Friends{Friends: map[string]map[string]struct{}{}},
		Money: storage.MoneyExchange{
			Owes:  map[string]storage.To{},
			Lends: map[string]storage.To{},
		},
		Groups:       storage.Groups{Groups: map[uuid.UUID]storage.Group{}},
		Participants: storage.Participants{Participants: map[uuid.UUID]map[string]struct{}{}},
		Participates: storage.Participates{Participates: map[string]map[uuid.UUID]struct{}{}},
		Server:       NewServer(),
	}
}

func Notify(a *App, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || !a.CheckPassword(username, password) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			r.Header.Set("string", username)
			f(w, r)
		}
	}
}

func (a *App) CreateUser(res http.ResponseWriter, req *http.Request) {
	username, password, _ := req.BasicAuth()

	err := a.Users.Create(string(username), password)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusCreated)
}

func (a *App) Login(res http.ResponseWriter, req *http.Request) {
	username, password, _ := req.BasicAuth()
	err := a.Users.CheckCredentials(username, password)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
	}
}

func (a *App) CheckPassword(username, password string) bool {
	realPassword := a.Users.GetPassword(username)
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

func (a *App) AddFriend(res http.ResponseWriter, req *http.Request) {
	headerContentType := req.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		http.Error(res, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	username := req.Header.Get("string")
	body, _ := ioutil.ReadAll(req.Body)

	var data map[string]string
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(res, "Error in Unmarshal", http.StatusInternalServerError)
		return
	}

	err := a.Friends.Add(username, data["friend"])
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusCreated)
}

func (a *App) ShowFriends(res http.ResponseWriter, req *http.Request) {
	username := req.Header.Get("string")
	friends := a.Friends.GetFriendsOf(username)
	marshal, _ := json.Marshal(friends)
	_, _ = res.Write(marshal)
}

func (a *App) AddDebt(res http.ResponseWriter, req *http.Request) {
	// TODO where to put?
	// Given: creditor, groupName, amount, reason
	// TODO get from body
	var (
		creditor  string
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

	// p.GetParticipants --> []string
	participants := a.Participants.GetParticipants(groupID)

	// g.AddDebt(creditor, groupID, participants, individualAmount, reason)
	a.Groups.AddDebt(creditor, groupID, participants, debt, reason)
}

func (a *App) CreateGroup(res http.ResponseWriter, req *http.Request) {
	headerContentType := req.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		http.Error(res, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	username := req.Header.Get("string")
	body, _ := ioutil.ReadAll(req.Body)

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(res, "Error in Unmarshal", http.StatusInternalServerError)
		return
	}

	groupName := data["name"].(string)
	participants := data["participants"].([]string)

	participants = append(participants, username)

	a.Groups.CreateGroup(groupName, participants)
	//a.Participants.Add()

	res.WriteHeader(http.StatusCreated)
}
