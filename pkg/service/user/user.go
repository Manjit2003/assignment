package user_service

import "github.com/Manjit2003/samespace/pkg/db"

var session = db.ScyllaSession

func AddUser(username, hashed_password string) (string, error) {
	session.Query("")

	return "", nil
}
