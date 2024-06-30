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
	queryGetTodos       = `SELECT id, title, description, status, created, updated FROM todos WHERE user_id = ?`
	queryGetSingleTodos = `SELECT id, title, description, status, created, updated FROM todos WHERE user_id = ? AND id = ? LIMIT 1`
	queryUpdateTodo     = `INSERT INTO todos (id, user_id, title, description, status, created, updated)
				VALUES (?, ?, ?, ?, ?, ?, toTimestamp(now()));`
	queryDeleteTodo = `DELETE FROM todos WHERE id = ? AND user_id = ?`
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
	query := db.ScyllaSession.Query(q, params...).PageSize(pageSize).PageState(pageState)

	iter := query.Iter()
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

func GetSingleUserTodo(userId, todoId string) (*model.TodoItem, error) {
	var todo model.TodoItem
	err := db.ScyllaSession.Query(queryGetSingleTodos, userId, todoId).Scan(
		&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.Created, &todo.Updated,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get todo: %v", err)
	}
	return &todo, nil
}

func UpdateUserTodo(item model.TodoItem, userId string) error {

	// the updates in scylla/cassandra cannot be made on primary keys.
	// status here is our primary key. so we'll have to remove that particular todo
	// an insert new one with same id but updated data :)

	if err := DeleteUserTodo(item.ID, userId); err != nil {
		return err
	}

	return db.ScyllaSession.Query(
		queryUpdateTodo,
		item.ID,
		userId,
		item.Title,
		item.Description,
		item.Status,
		item.Created,
	).Exec()
}

func PatchUserTodo(patch model.TodoPatch, todoId, userId string) error {
	item, err := GetSingleUserTodo(userId, todoId)

	if err != nil {
		return err
	}

	if patch.Title != nil {
		item.Title = *patch.Title
	}
	if patch.Description != nil {
		item.Description = *patch.Description
	}
	if patch.Status != nil {
		item.Status = *patch.Status
	}

	return UpdateUserTodo(*item, userId)
}

func DeleteUserTodo(todoId string, userId string) error {
	return db.ScyllaSession.Query(
		queryDeleteTodo,
		todoId,
		userId,
	).Exec()
}
