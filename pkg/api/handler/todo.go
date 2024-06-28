package handler

import (
	"context"
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/Manjit2003/samespace/pkg/api/middleware"
	"github.com/Manjit2003/samespace/pkg/model"
	todo_service "github.com/Manjit2003/samespace/pkg/service/todo"
	"github.com/Manjit2003/samespace/pkg/utils"
)

func GetUserIDFromContext(ctx context.Context) string {
	userID, _ := ctx.Value(middleware.UserKey).(string)
	return userID
}

// @Summary      Get user's TODOs in paginated form
// @Description  Returns a paginated list of TODO items for a user, with support for status filtering and pagination.
// @Tags         Todos
// @Produce      json
// @Param        status      query    string  false "TODO status (e.g., pending, completed)"
// @Param        page_state  query    string  false "Pagination state"
// @Param        page_size   query    int     false "Page size"
// @Success      200  {object}  utils.HTTPReponse  "Todos fetched"
// @Failure      400  {object}  utils.HTTPReponse  "Invalid request payload"
// @Failure      401  {object}  utils.HTTPReponse  "Unauthorized"
// @Failure      500  {object}  utils.HTTPReponse  "Internal server error"// @Security     ApiKeyAuth
// @Router       /todos [get]
// @Security     ApiKeyAuth
func HandleGetUserTodos(w http.ResponseWriter, r *http.Request) {

	userId := GetUserIDFromContext(r.Context())
	//status := r.URL.Query().Get("status")
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

	todos, nextPageState, err := todo_service.GetUserTodos(userId, pageSize, pageState, nil)

	if err != nil {
		utils.SendResponse(w, 500, utils.HTTPReponse{
			Error:   true,
			Message: "error getting todos",
		})
		return
	}

	utils.SendResponse(w, 200, utils.HTTPReponse{
		Error:   false,
		Message: "token refreshed",
		Data: map[string]interface{}{
			"data":      todos,
			"next_page": nextPageState,
		},
	})
}

// HandleAddUserTodo adds a new todo item for a user.
// @Summary      Add a new todo item
// @Description  Adds a new todo item for a user with the provided title and description. Requires user authentication.
// @Tags         Todos
// @Accept       json
// @Produce      json
// @Param        request  body  handler.HandleAddUserTodo.payload  true  "Todo item data"
// @Success      200  {object}  utils.HTTPReponse  "Todo item added successfully"
// @Failure      400  {object}  utils.HTTPReponse  "Invalid request payload"
// @Failure      401  {object}  utils.HTTPReponse  "Unauthorized"
// @Failure      500  {object}  utils.HTTPReponse  "Internal server error"
// @Router       /todos [post]
// @Security     ApiKeyAuth
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
