package user_service

import (
	"github.com/Manjit2003/samespace/pkg/db"
	"github.com/Manjit2003/samespace/pkg/model"
	"github.com/gocql/gocql"
)

const (
	queryGetUser = `SELECT id, username, hashed_password, created, updated FROM users
				WHERE username = ?;`
	queryAddUser = `INSERT INTO users (id, username, hashed_password, created, updated)
				VALUES (uuid(), ?, ?, toTimestamp(now()), toTimestamp(now()))`
)

func AddUser(username, hashedPassword string) error {
	return db.ScyllaSession.Query(queryAddUser, username, hashedPassword).Exec()
}

func GetUser(username string) (*model.User, error) {
	var user model.User
	if err := db.ScyllaSession.Query(queryGetUser, username).Consistency(gocql.One).Scan(
		&user.ID,
		&user.Username,
		&user.HashedPassword,
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
