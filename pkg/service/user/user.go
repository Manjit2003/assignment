package user_service

import (
	"github.com/Manjit2003/samespace/pkg/db"
	"github.com/Manjit2003/samespace/pkg/model"
	"github.com/gocql/gocql"
)

var session = db.ScyllaSession

const (
	queryGetUser = `SELECT id, username, created, updated FROM your_keyspace_name.users
              WHERE username = ? AND hashed_password = ?;`
	queryAddUser = `INSERT INTO users (id, username, hashed_password, created, updated)
              VALUES (uuid(), ?, ?, toTimestamp(now()), toTimestamp(now()))`
)

func AddUser(username, hashedPassword string) error {
	return session.Query(queryAddUser, username, hashedPassword).Exec()
}

func GetUser(username, hashedPassword string) (*model.User, error) {
	var user model.User
	if err := session.Query(queryGetUser, username, hashedPassword).Consistency(gocql.One).Scan(
		&user.ID,
		&user.Username,
		&user.Created,
		&user.Updated,
	); err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
