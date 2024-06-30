package config

var TestConfig = APIConfig{
	Web: webCfg{
		Host: "127.0.0.1",
		Port: 8080,
	},
	Database: dbCfg{
		Host:     "127.0.0.1",
		Port:     9042,
		Keyspace: "samespace_keyspace",
	},
	Auth: authCfg{
		JWTSecret:       "super_secret",
		TokenExpiration: 10,
	},
	Logging: logCfg{
		Level:    "info",
		Filename: "app.log",
	},
}
