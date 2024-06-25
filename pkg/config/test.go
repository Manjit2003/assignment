package config

import "github.com/Manjit2003/samespace/pkg/db"

var TestDBConfig = db.DBConfig{
	Hosts: "127.0.0.1",
	Port:  9042,
}
