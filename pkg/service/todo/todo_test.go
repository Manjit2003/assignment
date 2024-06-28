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

func TestAddTodo(t *testing.T) {

	db.InitDatabase(config.TestDBConfig)

	userId := createRandomUser()

	err := todo_service.AddUserTodo(model.TodoItem{
		Title:       "Hello world",
		Description: "Description",
		Status:      "pending",
	}, userId)

	if err != nil {
		t.Errorf("error creating a todo: %v", err)
	}
}

func TestGetTodos(t *testing.T) {

	db.InitDatabase(config.TestDBConfig)

	userId := createRandomUser()

	t.Run("create todo under user", func(t *testing.T) {
		err := todo_service.AddUserTodo(model.TodoItem{
			Title:       "Hello world",
			Description: "Description",
			Status:      "pending",
		}, userId)

		if err != nil {
			t.Errorf("error creating a todo: %v", err)
		}
	})

	t.Run("get user todos", func(t *testing.T) {
		todos, _, err := todo_service.GetUserTodos(userId, 10, []byte(""), nil)

		if err != nil {
			t.Errorf("error getting todos!: %v", err)
		}

		if len(todos) < 1 {
			t.Errorf("todos are blank")
		}
	})

	t.Run("get paginated list", func(t *testing.T) {
		for range 20 {
			todo_service.AddUserTodo(model.TodoItem{
				Title:       "Hello world",
				Description: "Description",
				Status:      "pending",
			}, userId)

		}

		todos, _, err := todo_service.GetUserTodos(userId, 10, []byte(""), nil)

		if err != nil {
			t.Errorf("error getting todos!: %v", err)
		}

		if len(todos) != 10 {
			t.Errorf("todos are more than requested :%d", len(todos))
		}

	})

}
