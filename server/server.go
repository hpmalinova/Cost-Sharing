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

// # Friends
func (a *App) AddFriend(res http.ResponseWriter, req *http.Request) {
	headerContentType := req.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		http.Error(res, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	username := req.Header.Get("Username")
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
	username := req.Header.Get("Username")
	friends := a.Friends.GetFriendsOf(username)
	marshal, _ := json.Marshal(friends)
	_, _ = res.Write(marshal)
}

// # Groups
func (a *App) CreateGroup(res http.ResponseWriter, req *http.Request) {
	headerContentType := req.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		http.Error(res, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	username := req.Header.Get("Username")
	body, _ := ioutil.ReadAll(req.Body)

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(res, "Error in Unmarshal", http.StatusInternalServerError)
		return
	}

	groupName := data["name"].(string)
	participantsInt := data["participants"].([]interface{})
	participants := make([]string, len(participantsInt))
	for i, _ := range participantsInt {
		participants[i] = participantsInt[i].(string)
	}

	participants = append(participants, username)

	groupID := a.Groups.CreateGroup(groupName, participants)

	a.Participants.Add(groupID, participants)

	for _, username := range participants {
		a.Participates.Add(username, groupID)
	}

	res.WriteHeader(http.StatusCreated)
}

func (a *App) ShowGroups(res http.ResponseWriter, req *http.Request) {
	username := req.Header.Get("Username")

	groupIDs := a.Participates.GetGroups(username)
	groupNames := a.Groups.GetGroupNames(groupIDs)

	marshal, _ := json.Marshal(groupNames)
	_, _ = res.Write(marshal)
}

// # Debts
func (a *App) AddDebtToFriend(res http.ResponseWriter, req *http.Request) {
	headerContentType := req.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		http.Error(res, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	username := req.Header.Get("Username")
	body, _ := ioutil.ReadAll(req.Body)

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(res, "Error in Unmarshal", http.StatusInternalServerError)
		return
	}

	friendName := data["friend"].(string)
	amount := int(data["amount"].(float64)) // todo?
	reason := data["reason"].(string)
	creditor := data["creditor"].(bool)

	if !a.Users.DoesExist(friendName) {
		msg := friendName + " does not exist!"
		http.Error(res, msg, http.StatusBadRequest)
		return
	}

	if creditor {
		a.Money.AddDebt(friendName, username, amount, reason)
	} else {
		a.Money.AddDebt(username, friendName, amount, reason)
	}

	res.WriteHeader(http.StatusCreated)
}

func (a *App) AddDebtToGroup(res http.ResponseWriter, req *http.Request) {
	headerContentType := req.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		http.Error(res, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	creditor := req.Header.Get("Username")
	body, _ := ioutil.ReadAll(req.Body)

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(res, "Error in Unmarshal", http.StatusInternalServerError)
		return
	}

	groupName := data["group"].(string)
	amount := int(data["amount"].(float64)) // todo?
	reason := data["reason"].(string)

	// Find GroupID of group with name <groupName>
	// In which groups participates the creditor
	groupIDs := a.Participates.GetGroups(creditor)

	groupID, err := a.Groups.FindGroupID(groupName, groupIDs)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	fAmount := math.Ceil(float64(amount) / float64(a.Participants.GetNumberOfParticipants(groupID)))
	amount = int(fAmount)

	participants := a.Participants.GetParticipants(groupID)

	a.Groups.AddDebt(creditor, groupID, participants, amount, reason)

	res.WriteHeader(http.StatusCreated)
}

func (a *App) ReturnDebt(res http.ResponseWriter, req *http.Request) {
	headerContentType := req.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		http.Error(res, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	debtor := req.Header.Get("Username")
	body, _ := ioutil.ReadAll(req.Body)

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(res, "Error in Unmarshal", http.StatusInternalServerError)
		return
	}

	friend := data["friend"].(string)
	amount := int(data["amount"].(float64)) // todo?
	groupName := data["groupName"].(string)

	groupIDs := a.Participates.GetGroups(debtor)

	groupID, err := a.Groups.FindGroupID(groupName, groupIDs)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the given friend is in the same group
	if !a.Participates.DoesParticipate(friend, groupID) {
		msg := friend + " is not part of " + groupName
		http.Error(res, msg, http.StatusBadRequest)
		return
	}

	a.Groups.ReturnDebt(debtor, groupID, friend, amount)

	res.WriteHeader(http.StatusCreated)
}

func (a *App) ShowOwed(res http.ResponseWriter, req *http.Request) {
	debtor := req.Header.Get("Username")

	owed := a.Money.GetOwed(debtor) // {to, amount, reason}
	marshal, _ := json.Marshal(owed)
	_, _ = res.Write(marshal)
}

func (a *App) ShowLent(res http.ResponseWriter, req *http.Request) {
	creditor := req.Header.Get("Username")

	lent := a.Money.GetLent(creditor) // {to, amount, reason}
	marshal, _ := json.Marshal(lent)
	_, _ = res.Write(marshal)
}
