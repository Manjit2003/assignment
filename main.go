package main

import (
	"fmt"

	"github.com/gocql/gocql"
)

func main() {

	cluster := gocql.NewCluster("127.0.0.1:9042")

	session, err := gocql.NewSession(*cluster)

	if err != nil {
		panic(err)
	}

	fmt.Println(session)
}
