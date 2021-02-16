package main

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

const (
	peter     = "peter"
	peterPass = "5678"
	george    = "george"
	lily      = "lily"
	maria     = "maria"
	juan = "juan"
	amount    = 20
	reason    = "food"
	presents = "presents"
	japan = "japan"
)

func (a *App) InitData(){
	users := storage.Users{Users: map[string]storage.User{}}
	_ = users.Create(peter, peterPass)
	_ = users.Create(george, "7890")
	_ = users.Create(lily, "1234")
	_ = users.Create(maria, "0000")
	_ = users.Create(juan, "8888")
	a.Users = users

	friends := storage.Friends{Friends: map[string]map[string]struct{}{}}
	_ = friends.Add(peter, george)
	_ = friends.Add(peter, lily)
	a.Friends = friends

	moneyExchange := storage.MoneyExchange{Owes: map[string]storage.To{}, Lends: map[string]storage.To{}}
	moneyExchange.AddUser(peter)
	moneyExchange.AddUser(george)
	moneyExchange.AddUser(lily)
	moneyExchange.AddDebt(peter, george, amount, reason)
	moneyExchange.AddDebt(lily, peter, amount, reason)
	a.Money = moneyExchange

	groups := storage.Groups{Groups: map[uuid.UUID]storage.Group{}}
	p1 := []string{peter, lily}
	id1 := groups.CreateGroup(presents, p1)
	groups.AddDebt(peter, id1, p1, amount, reason)
	p2 := []string{peter, george, maria}
	id2 := groups.CreateGroup(japan, p2)
	groups.AddDebt(peter, id2, p2, amount, reason)
	a.Groups = groups

	participants := storage.Participants{Participants: map[uuid.UUID]map[string]struct{}{}}
	participants.Add(id1, p1)
	participants.Add(id2, p2)
	a.Participants = participants

	participates := storage.Participates{map[string]map[uuid.UUID]struct{}{}}
	participates.Add(peter, id1)
	participates.Add(lily, id1)
	participates.Add(peter, id2)
	participates.Add(george, id2)
	participates.Add(maria, id2)
	a.Participates = participates
}

func InitApp(withInitData bool) App {
	a := App{
		Users:        storage.Users{Users: map[string]storage.User{}},
		Friends:      storage.Friends{Friends: map[string]map[string]struct{}{}},
		Money:        storage.MoneyExchange{Owes: map[string]storage.To{}, Lends: map[string]storage.To{}},
		Groups:       storage.Groups{Groups: map[uuid.UUID]storage.Group{}},
		Participants: storage.Participants{Participants: map[uuid.UUID]map[string]struct{}{}},
		Participates: storage.Participates{Participates: map[string]map[uuid.UUID]struct{}{}},
		Server:       NewServer(),
	}
	if withInitData{
		a.InitData()
	}
	return a
}

// Notify works as a middleware. If the given credentials are not valid,
// returns http.StatusUnauthorized. Else, writes a Header(Key: "Username")
func Notify(a *App, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok || !a.checkPassword(username, password) {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			r.Header.Set("Username", username)
			f(w, r)
		}
	}
}

// # Users

// CreateAccount returns http.StatusBadRequest if a user with the same name already exists
// and http.StatusCreated if the new user is created successfully
func (a *App) CreateAccount(res http.ResponseWriter, req *http.Request) {
	username, password, _ := req.BasicAuth()

	err := a.Users.Create(username, password)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusCreated)
}

// Login returns http.StatusUnauthorized if the credentials are not valid
// and http.StatusOK if they are
func (a *App) Login(res http.ResponseWriter, req *http.Request) {
	username, password, _ := req.BasicAuth()
	err := a.Users.CheckCredentials(username, password)
	if err != nil {
		http.Error(res, err.Error(), http.StatusUnauthorized)
	}
}

// ShowUsers writes to the response.Body:
// []string that contains all of the usernames in the storage
func (a *App) ShowUsers(res http.ResponseWriter, req *http.Request) {
	marshal, _ := json.Marshal(a.Users.GetUsernames())
	_, _ = res.Write(marshal)
}

// # Friends

// AddFriend expects map[string]string{"friend": friendName} from the request.
// If the content type is not "application/json", returns http.StatusUnsupportedMediaType.
// If the client is already friends with that "friend", returns http.StatusBadRequest.
// If the request is successful, returns http.StatusCreated
func (a *App) AddFriend(res http.ResponseWriter, req *http.Request) {
	headerContentType := req.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		http.Error(res, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	username := req.Header.Get("Username")
	body, _ := ioutil.ReadAll(req.Body)

	var data map[string]string
	_ = json.Unmarshal(body, &data)

	err := a.Friends.Add(username, data["friend"])
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.WriteHeader(http.StatusCreated)
}

// ShowFriends writes to the response.Body
// []string that contains all of the client friends
func (a *App) ShowFriends(res http.ResponseWriter, req *http.Request) {
	username := req.Header.Get("Username")
	friends := a.Friends.GetFriendsOf(username)
	marshal, _ := json.Marshal(friends)
	_, _ = res.Write(marshal)
}

// AddDebtToFriend receives map[string]interface{} that contains
// {"friend": string, "amount": int, "reason": string, "creditor":bool},
// where creditor is True, if the client wants to give money to his friend,
// and False , if the clients takes money from his friend.
// If the content type is not "application/json", returns http.StatusUnsupportedMediaType.
// If there is a problem with the Unmarshal, returns http.StatusInternalServerError
// If the client and the username are not friends, returns http.StatusBadRequest.
// If the request is successful, returns http.StatusCreated.
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
	amount := int(data["amount"].(float64))
	reason := data["reason"].(string)
	creditor := data["creditor"].(bool)

	if !a.Friends.AreFriends(username, friendName) {
		msg := friendName + " is not your friend!"
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

// ShowDebts writes to the response.Body:
// DebtC{To,Amount,Reason} for each debt the client has
func (a *App) ShowDebts(res http.ResponseWriter, req *http.Request) {
	debtor := req.Header.Get("Username")

	owed := a.Money.GetOwed(debtor)
	marshal, _ := json.Marshal(owed)
	_, _ = res.Write(marshal)
}

// ShowLoans writes to the response.Body:
// DebtC{To,Amount,Reason} for each loan the client has given
func (a *App) ShowLoans(res http.ResponseWriter, req *http.Request) {
	creditor := req.Header.Get("Username")

	lent := a.Money.GetLent(creditor) // {to, amount, reason}
	marshal, _ := json.Marshal(lent)
	_, _ = res.Write(marshal)
}


// # Groups

// CreateGroup receives map[string]interface{}{"name": groupName, "participants": []usernames}
// If the content type is not "application/json", returns http.StatusUnsupportedMediaType.
// If there is a problem with the Unmarshal, returns http.StatusInternalServerError.
// If the request is successful, returns http.StatusCreated.
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

// ShowGroups writes to the response.Body
// []string that contains all of the groupNames that the client participates in
func (a *App) ShowGroups(res http.ResponseWriter, req *http.Request) {
	username := req.Header.Get("Username")

	groupIDs := a.Participates.GetGroups(username)
	groupNames := a.Groups.GetGroupNames(groupIDs)

	marshal, _ := json.Marshal(groupNames)
	_, _ = res.Write(marshal)
}

// AddDebtToGroup receives map[string]interface{} with keys: {"group", "amount", "reason"}
// If the content type is not "application/json", returns http.StatusUnsupportedMediaType.
// If there is a problem with the Unmarshal, returns http.StatusInternalServerError.
// If the client does not participate in group with that name, return http.StatusBadRequest.
// Divides the amount of bebt into equal parts (int of ceil(result)) among all of the participants
// If the request is successful, returns http.StatusCreated.
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
	amount := int(data["amount"].(float64))
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

// If there is a problem with the Unmarshal, returns http.StatusInternalServerError
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

// ShowDebtsToGroups writes to the response.Body:
// map[groupName][]DebtC{To,Amount,Reason} for each debt the client has to a group
func (a *App) ShowDebtsToGroups(res http.ResponseWriter, req *http.Request) {
	debtor := req.Header.Get("Username")
	groupIDs := a.Participates.GetGroups(debtor)

	owed := a.Groups.GetOwed(debtor, groupIDs)
	marshal, _ := json.Marshal(owed)
	_, _ = res.Write(marshal)
}

// ShowLoansToGroups writes to the response.Body:
// map[groupName][]DebtC{To,Amount,Reason} for each loan the client has given to a group
func (a *App) ShowLoansToGroups(res http.ResponseWriter, req *http.Request) {
	creditor := req.Header.Get("Username")
	groupIDs := a.Participates.GetGroups(creditor)

	lent := a.Groups.GetLent(creditor, groupIDs)
	marshal, _ := json.Marshal(lent)
	_, _ = res.Write(marshal)
}

// # Utils
func (a *App) checkPassword(username, password string) bool {
	realPassword := a.Users.GetPassword(username)
	return checkPasswordHash(password, realPassword)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
