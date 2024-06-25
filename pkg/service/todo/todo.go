package service

import (
	"fmt"

	"github.com/Manjit2003/samespace/pkg/db"
	"github.com/Manjit2003/samespace/pkg/model"
)

const (
	queryAddTodo = `INSERT INTO todos (id, user_id, title, description, status, created, updated)
				VALUES (uuid(), ?, ?, ?, ? toTimestamp(now()), toTimestamp(now()))`
)

func AddUserTodo(item model.TodoItem, userId string) error {
	return db.ScyllaSession.Query(
		queryAddTodo,
		item.UserID,
		item.Title,
		item.Description,
		item.Status,
	).Exec()
}

func GetUserTodos(userId string, pageSize int, pageState string, status *string) ([]model.TodoItem, string, error) {
	q := "SELECT id, title, description, status, created, updated FROM todo_items WHERE user_id = ?"
	params := []interface{}{userId}
	if status != nil {
		q += " AND status = ?"
		params = append(params, *status)
	}
	var todos []model.TodoItem
	query := db.ScyllaSession.Query(q, params...).PageSize(pageSize)
	if pageState != "" {
		query = query.PageState([]byte(pageState))
	}
	iter := query.Iter()
	var todo model.TodoItem
	for iter.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.Created, &todo.Updated) {
		todos = append(todos, todo)
	}
	newPageState := string(iter.PageState())
	if err := iter.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to query todos: %v", err)
	}
	return todos, newPageState, nil
}
