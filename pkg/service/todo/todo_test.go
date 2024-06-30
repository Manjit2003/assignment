package todo_service_test

import (
	"testing"

	"github.com/Manjit2003/samespace/pkg/config"
	"github.com/Manjit2003/samespace/pkg/db"
	"github.com/Manjit2003/samespace/pkg/model"
	todo_service "github.com/Manjit2003/samespace/pkg/service/todo"
	user_service "github.com/Manjit2003/samespace/pkg/service/user"
	"github.com/Manjit2003/samespace/pkg/utils"
)

func createRandomUser() string {
	username, pass := utils.GenerateRandomCreds()
	user_service.AddUser(username, pass)
	user, _ := user_service.GetUser(username)
	return user.ID
}

func cleanDatabase() {
	db.ScyllaSession.Query(`TRUNCATE todos`).Exec()
	db.ScyllaSession.Query(`TRUNCATE users`).Exec()
}

func TestTodoService(t *testing.T) {
	db.InitDatabase(&config.TestConfig)
	defer cleanDatabase()

	userId := createRandomUser()

	t.Run("create todo for user", func(t *testing.T) {
		err := todo_service.AddUserTodo(model.TodoItem{
			Title:       "Hello world",
			Description: "Description",
			Status:      "pending",
		}, userId)

		if err != nil {
			t.Fatalf("error creating a todo: %v", err)
		}
	})

	t.Run("get user todos", func(t *testing.T) {
		todos, _, err := todo_service.GetUserTodos(userId, 10, nil, nil, "")

		if err != nil {
			t.Fatalf("error getting todos: %v", err)
		}

		if len(todos) < 1 {
			t.Fatalf("expected at least one todo, got %d", len(todos))
		}
	})

	t.Run("get paginated list of todos", func(t *testing.T) {
		for i := 0; i < 20; i++ {
			todo_service.AddUserTodo(model.TodoItem{
				Title:       "Hello world",
				Description: "Description",
				Status:      "pending",
			}, userId)
		}

		status := "pending"
		todos, nextPageState, err := todo_service.GetUserTodos(userId, 10, nil, &status, "")

		if err != nil {
			t.Fatalf("error getting todos: %v", err)
		}

		if len(todos) != 10 {
			t.Fatalf("expected 10 todos, got %d", len(todos))
		}

		if nextPageState == "" {
			t.Fatalf("expected non-empty nextPageState for pagination")
		}
	})

	t.Run("get single todo of user", func(t *testing.T) {
		todos, _, err := todo_service.GetUserTodos(userId, 1, nil, nil, "")

		if err != nil {
			t.Fatalf("error getting todos: %v", err)
		}

		if len(todos) < 1 {
			t.Fatalf("expected at least one todo to fetch a single todo")
		}

		todoId := todos[0].ID
		singleTodo, err := todo_service.GetSingleUserTodo(userId, todoId)

		if err != nil {
			t.Fatalf("error getting single todo: %v", err)
		}

		if singleTodo.ID != todoId {
			t.Fatalf("returned wrong todo: expected %s, received %s", todoId, singleTodo.ID)
		}
	})

	t.Run("update todo of user", func(t *testing.T) {
		todos, _, err := todo_service.GetUserTodos(userId, 1, nil, nil, "")

		if err != nil {
			t.Fatalf("error getting todos: %v", err)
		}

		if len(todos) < 1 {
			t.Fatalf("expected at least one todo to update")
		}

		todo := todos[0]
		todo.Title = "Updated title"
		todo.Description = "Updated description"
		todo.Status = "completed"

		err = todo_service.UpdateUserTodo(todo, userId)

		if err != nil {
			t.Fatalf("error updating todo: %v", err)
		}

		updatedTodo, err := todo_service.GetSingleUserTodo(userId, todo.ID)

		if err != nil {
			t.Fatalf("error getting updated todo: %v", err)
		}

		if updatedTodo.Title != "Updated title" || updatedTodo.Description != "Updated description" || updatedTodo.Status != "completed" {
			t.Fatalf("todo not updated correctly: %+v", updatedTodo)
		}
	})

	t.Run("patch todo of user", func(t *testing.T) {
		todos, _, err := todo_service.GetUserTodos(userId, 1, nil, nil, "")

		if err != nil {
			t.Fatalf("error getting todos: %v", err)
		}

		if len(todos) < 1 {
			t.Fatalf("expected at least one todo to patch")
		}

		todo := todos[0]
		patch := model.TodoPatch{
			Title:       utils.StringPtr("Patched title"),
			Description: utils.StringPtr("Patched description"),
			Status:      utils.StringPtr("in-progress"),
		}

		err = todo_service.PatchUserTodo(patch, todo.ID, userId)

		if err != nil {
			t.Fatalf("error patching todo: %v", err)
		}

		patchedTodo, err := todo_service.GetSingleUserTodo(userId, todo.ID)

		if err != nil {
			t.Fatalf("error getting patched todo: %v", err)
		}

		if patchedTodo.Title != "Patched title" || patchedTodo.Description != "Patched description" || patchedTodo.Status != "in-progress" {
			t.Fatalf("todo not patched correctly: %+v", patchedTodo)
		}
	})

	t.Run("delete todo of user", func(t *testing.T) {
		todos, _, err := todo_service.GetUserTodos(userId, 1, nil, nil, "created.asc")

		if err != nil {
			t.Fatalf("error getting todos: %v", err)
		}

		if len(todos) < 1 {
			t.Fatalf("expected at least one todo to delete")
		}

		todo := todos[0]

		err = todo_service.DeleteUserTodo(todo.ID, userId)

		if err != nil {
			t.Fatalf("error deleting todo: %v", err)
		}

		_, err = todo_service.GetSingleUserTodo(userId, todo.ID)

		if err == nil {
			t.Fatalf("expected error when getting deleted todo, got none")
		}
	})
}
