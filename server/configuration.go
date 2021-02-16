package main

const (
	UrlLogin             = Protocol + Host + Port + PathToLogin
	UrlCreateAccount     = Protocol + Host + Port + PathToCreateAccount
	UrlShowUsers         = Protocol + Host + Port + PathToShowUsers
	UrlAddFriend         = Protocol + Host + Port + PathToAddFriend
	UrlShowFriends       = Protocol + Host + Port + PathToShowFriends
	UrlCreateGroup       = Protocol + Host + Port + PathToCreateGroup
	UrlShowGroups        = Protocol + Host + Port + PathToShowGroups
	UrlAddDebt           = Protocol + Host + Port + PathToAddDebt
	UrlAddDebtToGroup    = Protocol + Host + Port + PathToAddDebtToGroup
	UrlReturnDebt        = Protocol + Host + Port + PathToReturnDebt
	UrlShowDebts         = Protocol + Host + Port + PathToShowDebts
	UrlShowLoans         = Protocol + Host + Port + PathToShowLoans
	UrlShowDebtsToGroups = Protocol + Host + Port + PathToShowDebtsToGroups
	UrlShowLoansToGroups = Protocol + Host + Port + PathToShowLoansToGroups
)

const (
	Protocol                = "http://"
	Host                    = "localhost"
	Port                    = ":8080"
	PathToCreateAccount     = "/costSharing/createAccount"
	PathToLogin             = "/costSharing/login"
	PathToShowUsers         = "/costSharing/home/showUsers"
	PathToAddFriend         = "/costSharing/home/addFriend"
	PathToShowFriends       = "/costSharing/home/showFriends"
	PathToCreateGroup       = "/costSharing/home/createGroup"
	PathToShowGroups        = "/costSharing/home/showGroups"
	PathToAddDebt           = "/costSharing/home/addDebt"
	PathToShowDebts         = "/costSharing/home/showDebts"
	PathToShowLoans         = "/costSharing/home/showLoans"
	PathToAddDebtToGroup    = "/costSharing/home/addDebtToGroup"
	PathToReturnDebt        = "/costSharing/home/returnDebt"
	PathToShowDebtsToGroups = "/costSharing/home/showDebtsToGroups"
	PathToShowLoansToGroups = "/costSharing/home/showLoansToGroups"
)
