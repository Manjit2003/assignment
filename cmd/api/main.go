package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/Manjit2003/samespace/pkg/api/server"
	"github.com/Manjit2003/samespace/pkg/config"
	"github.com/Manjit2003/samespace/pkg/db"
	"github.com/Manjit2003/samespace/pkg/utils"

	// http-swagger middleware
	_ "github.com/Manjit2003/samespace/docs"
)

// @title						Samespace Assignment
// @version					1.0
// @description				This is my submission to the assignment as requested by Samespace.
// @contact.name				Manjeet Pardeshi
// @contact.email				manjit2003@proton.me
// @BasePath					/api/v1
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @externalDocs.description	Github Repo
// @externalDocs.url			http://github.com/Manjit2003/assignment
// @tags Tag1, Tag2, Tag3
// @defaultExpandDepth 1
func main() {

	log := utils.GetChildLogger("entrypoint")

	config.LoadConfig("/root/config.yaml")
	time.Sleep(5 * time.Second)
	db.InitDatabase(config.Config)

	srv := server.MakeServer()

	var wait time.Duration = time.Second * 15

	listen := func() error {
		log.Info("server started successfully")
		return srv.ListenAndServe()
	}

	go func() {
		if err := listen(); err != nil {
			log.Error(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(ctx)
	log.Info("gracefully shutting down")
	os.Exit(0)
}
