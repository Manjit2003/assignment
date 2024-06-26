package handler

import (
	"net/http"

	"github.com/Manjit2003/samespace/pkg/model"
	"github.com/Manjit2003/samespace/pkg/utils"
)

func HandleUserLogin(w http.ResponseWriter, r *http.Request) {

	var payload model.AuthLoginData
	err := utils.ReadBody(w, r, &payload)
	if err != nil {
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Todo item created successfully"))

}
