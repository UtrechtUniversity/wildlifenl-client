package wildlifenl

import "time"

type Credential struct {
	UserID    string    `json:"userID"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	LastLogin time.Time `json:"lastLogon"`
	Scopes    []string  `json:"scopes"`
}
