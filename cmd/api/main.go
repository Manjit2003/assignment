package main

import (
	"fmt"

	"github.com/Manjit2003/samespace/pkg/db"
)

func main() {

	db.InitDatabase(db.DBConfig{
		Hosts: "127.0.0.1",
		Port:  9042,
	})

	fmt.Println(db.ScyllaSession)

}
