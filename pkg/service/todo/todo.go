package todo_service

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/Manjit2003/samespace/pkg/db"
	"github.com/Manjit2003/samespace/pkg/model"
)

const (
	queryAddTodo = `INSERT INTO todos (id, user_id, title, description, status, created, updated)
				VALUES (uuid(), ?, ?, ?, ?, toTimestamp(now()), toTimestamp(now()))`
	queryGetTodos = `SELECT id, title, description, status, created, updated FROM todos WHERE user_id = ?`
)

func AddUserTodo(item model.TodoItem, userId string) error {
	return db.ScyllaSession.Query(
		queryAddTodo,
		userId,
		item.Title,
		item.Description,
		item.Status,
	).Exec()
}

func GetUserTodos(userId string, pageSize int, pageState []byte, status *string) ([]model.TodoItem, string, error) {
	q := strings.Clone(queryGetTodos)
	params := []interface{}{userId}
	if status != nil {
		q += " AND status = ?"
		params = append(params, *status)
	}
	var todos []model.TodoItem
	qu := db.ScyllaSession.Query(q, params...).PageSize(10).PageState(nil)
	iter := qu.Iter()
	var todo model.TodoItem
	for iter.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.Created, &todo.Updated) {
		todos = append(todos, todo)
	}
	newPageState := base64.StdEncoding.EncodeToString(iter.PageState())
	if err := iter.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to query todos: %v", err)
	}
	return todos, newPageState, nil
}
