package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Manjit2003/samespace/pkg/api/middleware"
	"github.com/Manjit2003/samespace/pkg/model"
	todo_service "github.com/Manjit2003/samespace/pkg/service/todo"
	"github.com/Manjit2003/samespace/pkg/utils"
	"github.com/gorilla/mux"
)

func GetUserIDFromContext(ctx context.Context) string {
	userID, _ := ctx.Value(middleware.UserKey).(string)
	return userID
}

// HandleGetUserTodos handles the GET request to retrieve todo items
// @Summary		Get user's TODOs in paginated form
// @Description	Returns a paginated list of TODO items for a user, with support for status filtering and pagination.
// @Tags			Todos
// @Produce		json
// @Param			status		query		string				false	"TODO status (e.g., pending, completed)"
// @Param			sort		query		string				false	"The field to sort. (e.g., created.asc, updated.asc)"
// @Param			page_state	query		string				false	"Pagination state"
// @Param			page_size	query		int					false	"Page size"
// @Success		200			{object}	utils.HTTPReponse	"Todos fetched"
// @Failure		400			{object}	utils.HTTPReponse	"Invalid request payload"
// @Failure		401			{object}	utils.HTTPReponse	"Unauthorized"
// @Failure		500			{object}	utils.HTTPReponse	"Internal server error"
// @Security		ApiKeyAuth
// @Router			/todos [get]
func HandleGetUserTodos(w http.ResponseWriter, r *http.Request) {
	userId := GetUserIDFromContext(r.Context())
	status := r.URL.Query().Get("status")
	sort := r.URL.Query().Get("sort")
	pageStateStr := r.URL.Query().Get("page_state")
	pageSizeStr := r.URL.Query().Get("page_size")

	pageSize := 10 // setting it default to 10 for now :)
	var pageState []byte

	if pageStateStr != "" {
		ps, err := base64.StdEncoding.DecodeString(pageStateStr)
		if err == nil {
			pageState = ps
		}
	}

	if pageSizeStr != "" {
		ps, err := strconv.Atoi(pageSizeStr)
		if err == nil {
			pageSize = ps
		}
	}

	var todos []model.TodoItem
	var nextPageState string
	var err error
	if status != "" {
		todos, nextPageState, err = todo_service.GetUserTodos(userId, pageSize, pageState, &status, sort)
	} else {
		todos, nextPageState, err = todo_service.GetUserTodos(userId, pageSize, pageState, nil, sort)
	}

	if err != nil {
		utils.SendResponse(w, 500, utils.HTTPReponse{
			Error:   true,
			Message: "error getting todos",
		})
		return
	}

	utils.SendResponse(w, 200, utils.HTTPReponse{
		Error:   false,
		Message: "todos fetched",
		Data: map[string]interface{}{
			"data":      todos,
			"next_page": nextPageState,
		},
	})
}

// HandleAddUserTodo adds a new todo item for a user
// @Summary		Add a new todo item
// @Description	Adds a new todo item for a user with the provided title and description. Requires user authentication.
// @Tags			Todos
// @Accept			json
// @Produce		json
// @Param			request	body		handler.HandleAddUserTodo.payload	true	"Todo item data"
// @Success		200		{object}	utils.HTTPReponse					"Todo item added successfully"
// @Failure		400		{object}	utils.HTTPReponse					"Invalid request payload"
// @Failure		401		{object}	utils.HTTPReponse					"Unauthorized"
// @Failure		500		{object}	utils.HTTPReponse					"Internal server error"
// @Router			/todos [post]
// @Security		ApiKeyAuth
func HandleAddUserTodo(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	var data payload

	err := utils.ReadBody(w, r, &data)
	if err != nil {
		utils.SendResponse(w, 400, utils.HTTPReponse{
			Error:   true,
			Message: "invalid request payload",
		})
		return
	}

	userId := GetUserIDFromContext(r.Context())

	err = todo_service.AddUserTodo(model.TodoItem{
		Title:       data.Title,
		Description: data.Description,
		Status:      "pending",
	}, userId)

	if err != nil {
		utils.SendResponse(w, 500, utils.HTTPReponse{
			Error:   true,
			Message: "error adding new todo item",
		})
		return
	}

	utils.SendResponse(w, 200, utils.HTTPReponse{
		Error:   false,
		Message: "todo added",
	})
}

// HandleUpdateUserTodo updates a todo item for a user
// @Summary		Update a todo item
// @Description	Updates a todo item for a user with the provided title, description, and status. Requires user authentication.
// @Tags			Todos
// @Accept			json
// @Produce		json
// @Param			id		path		string							true	"Todo ID"
// @Param			request	body		handler.HandleUpdateUserTodo.payload	true	"Todo item data"
// @Success		200		{object}	utils.HTTPReponse				"Todo item updated successfully"
// @Failure		400		{object}	utils.HTTPReponse				"Invalid request payload"
// @Failure		401		{object}	utils.HTTPReponse				"Unauthorized"
// @Failure		500		{object}	utils.HTTPReponse				"Internal server error"
// @Router			/todos/{id} [put]
// @Security		ApiKeyAuth
func HandleUpdateUserTodo(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
	}
	var data payload

	err := utils.ReadBody(w, r, &data)
	if err != nil {
		utils.SendResponse(w, 400, utils.HTTPReponse{
			Error:   true,
			Message: "invalid request payload",
		})
		return
	}

	vars := mux.Vars(r)
	itemID := vars["id"]
	userId := GetUserIDFromContext(r.Context())

	err = todo_service.UpdateUserTodo(model.TodoItem{
		ID:          itemID,
		Title:       data.Title,
		Description: data.Description,
		Status:      data.Status,
	}, userId)

	if err != nil {
		utils.SendResponse(w, 500, utils.HTTPReponse{
			Error:   true,
			Message: "error updating todo item",
		})
		return
	}

	utils.SendResponse(w, 200, utils.HTTPReponse{
		Error:   false,
		Message: "todo updated",
	})
}

// HandlePatchUserTodo partially updates a todo item for a user
// @Summary		Partially update a todo item
// @Description	Partially updates a todo item for a user with the provided fields. Requires user authentication.
// @Tags			Todos
// @Accept			json
// @Produce		json
// @Param			id		path		string							true	"Todo ID"
// @Param			request	body		handler.HandlePatchUserTodo.payload	true	"Partial todo item data"
// @Success		200		{object}	utils.HTTPReponse				"Todo item updated successfully"
// @Failure		400		{object}	utils.HTTPReponse				"Invalid request payload"
// @Failure		401		{object}	utils.HTTPReponse				"Unauthorized"
// @Failure		500		{object}	utils.HTTPReponse				"Internal server error"
// @Router			/todos/{id} [patch]
// @Security		ApiKeyAuth
func HandlePatchUserTodo(w http.ResponseWriter, r *http.Request) {
	type payload map[string]interface{}
	var updates payload
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		utils.SendResponse(w, 400, utils.HTTPReponse{
			Error:   true,
			Message: "invalid request payload",
		})
		return
	}

	userId := GetUserIDFromContext(r.Context())

	vars := mux.Vars(r)

	// Get the existing todo item
	todo, err := todo_service.GetSingleUserTodo(userId, vars["id"])
	if err != nil {
		fmt.Println(err)
		utils.SendResponse(w, 404, utils.HTTPReponse{
			Error:   true,
			Message: "todo not found",
		})
		return
	}

	// Update the fields based on the PATCH request
	if title, ok := updates["title"]; ok {
		todo.Title = title.(string)
	}
	if description, ok := updates["description"]; ok {
		todo.Description = description.(string)
	}
	if status, ok := updates["status"]; ok {
		todo.Status = status.(string)
	}

	if err := todo_service.UpdateUserTodo(*todo, userId); err != nil {
		fmt.Println(err)
		utils.SendResponse(w, 500, utils.HTTPReponse{
			Error:   true,
			Message: "error updating todo item",
		})
		return
	}

	utils.SendResponse(w, 200, utils.HTTPReponse{
		Error:   false,
		Message: "todo updated",
	})
}

// HandleDeleteUserTodo deletes a todo item for a user
// @Summary		Delete a todo item
// @Description	Deletes a todo item for a user by its ID. Requires user authentication.
// @Tags			Todos
// @Produce		json
// @Param			id		path		string							true	"Todo ID"
// @Success		200		{object}	utils.HTTPReponse				"Todo item deleted successfully"
// @Failure		400		{object}	utils.HTTPReponse				"Invalid request payload"
// @Failure		401		{object}	utils.HTTPReponse				"Unauthorized"
// @Failure		500		{object}	utils.HTTPReponse				"Internal server error"
// @Router			/todos/{id} [delete]
// @Security		ApiKeyAuth
func HandleDeleteUserTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["id"]
	userId := GetUserIDFromContext(r.Context())

	if err := todo_service.DeleteUserTodo(todoId, userId); err != nil {
		utils.SendResponse(w, 500, utils.HTTPReponse{
			Error:   true,
			Message: "error deleting todo item",
		})
		return
	}

	utils.SendResponse(w, 200, utils.HTTPReponse{
		Error:   false,
		Message: "todo deleted",
	})
}
