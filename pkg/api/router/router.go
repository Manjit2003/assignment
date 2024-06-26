package router

import (
	"net/http"

	"github.com/Manjit2003/samespace/pkg/api/middleware"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func handler(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hellow world"))
}

func MakeRouter() *mux.Router {

	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	r.Use(middleware.LoggingMiddleware)

	return r

}
