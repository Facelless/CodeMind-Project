package match

import "miservicegolang/core/domain/user"

type MatchEvent struct {
	Type    string
	Message string
	User    user.User
}
