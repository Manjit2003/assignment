package router

import (
	"net/http"

	"github.com/Manjit2003/samespace/pkg/api/handler"
	"github.com/Manjit2003/samespace/pkg/api/middleware"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func MakeRouter() *mux.Router {

	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	v1 := r.PathPrefix("/api/v1").Subrouter()
	{
		v1.Use(middleware.LoggingMiddleware)

		authRouter := v1.PathPrefix("/auth").Subrouter()
		authRouter.HandleFunc("/login", handler.HandleUserLogin).Methods("POST")
		authRouter.HandleFunc("/register", handler.HandleUserRegister).Methods("POST")
		authRouter.HandleFunc("/refresh", handler.HandleGetAccessToken).Methods("POST")
	}

	return r
}
