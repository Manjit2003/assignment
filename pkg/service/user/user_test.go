package user_service_test

import (
	"testing"

	"github.com/Manjit2003/samespace/pkg/config"
	"github.com/Manjit2003/samespace/pkg/db"
	user_service "github.com/Manjit2003/samespace/pkg/service/user"
)

func TestCreateUser(t *testing.T) {
	db.InitDatabase(&config.TestConfig)

	if err := user_service.AddUser("Manjit2003", "my_sample_hashed_passwords"); err != nil {
		t.Errorf("err creating user: %v", err)
	}
}

func TestGetUser(t *testing.T) {
	db.InitDatabase(&config.TestConfig)

	if err := user_service.AddUser("Manjit2004", "my_sample_hashed_passwords"); err != nil {
		t.Errorf("err creating user: %v", err)
	}

	user, err := user_service.GetUser("Manjit2004")

	if err != nil {

		t.Errorf("error getting user from database: %v", err)
	}

	if user == nil {
		t.Errorf("user returned as null")
	}
}
