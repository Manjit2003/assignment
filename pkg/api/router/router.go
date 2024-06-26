package router

import "github.com/gorilla/mux"

func MakeRouter() *mux.Router {

	return mux.NewRouter()

}
