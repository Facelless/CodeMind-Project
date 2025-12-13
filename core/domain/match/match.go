package match

import (
	"miservicegolang/core/domain/user"
)

type Match struct {
	Id        string
	Type      string `json:"type"`
	User1     *user.User
	User2     *user.User
	Challenge string
}
