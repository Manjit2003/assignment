package server

import (
	"net/http"
	"time"

	"github.com/Manjit2003/samespace/pkg/api/router"
)

func MakeServer() *http.Server {
	r := router.MakeRouter()
	return &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}
}
