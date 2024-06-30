.
├── Makefile
├── bin
│   └── api
├── cmd
│   └── api
│       └── main.go
├── config.example.yaml
├── config.yaml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── main.go
└── pkg
    ├── api
    │   ├── handler
    │   │   ├── auth.go
    │   │   └── todo.go
    │   ├── middleware
    │   │   ├── auth.go
    │   │   └── log.go
    │   ├── router
    │   │   └── router.go
    │   └── server
    │       └── server.go
    ├── config
    │   ├── config.go
    │   └── test.go
    ├── db
    │   └── scylla.go
    ├── model
    │   ├── auth.go
    │   ├── todo.go
    │   └── user.go
    ├── service
    │   ├── auth
    │   │   ├── auth.go
    │   │   ├── auth_test.go
    │   │   └── token.go
    │   ├── todo
    │   │   ├── todo.go
    │   │   └── todo_test.go
    │   └── user
    │       ├── user.go
    │       └── user_test.go
    └── utils
        ├── http.go
        ├── logger.go
        └── test.go

19 directories, 33 files