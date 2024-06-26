package router

import (
	"net/http"

	"github.com/Manjit2003/samespace/pkg/api/middleware"
	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	return
}

func MakeRouter() *mux.Router {

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.AuthMiddleware)

	r.HandleFunc("/", handler)

	return r

}
