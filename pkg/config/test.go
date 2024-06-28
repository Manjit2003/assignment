package config

import "github.com/Manjit2003/samespace/pkg/db"

var TestDBConfig = db.DBConfig{
	Hosts: "scylla",
	Port:  9042,
}
