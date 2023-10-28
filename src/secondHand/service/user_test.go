package service

import (
	"errors"
	"secondHand/backend"
	"secondHand/model"
	"secondHand/util"
	"testing"
)

func TestCheckUserMatch(t *testing.T) {
	backend.InitPostgreSQLBackend()
	result, err := CheckUser("alice@alice", "alice")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result == false {
		t.Errorf("Result should be true, but get false")
	}
}

func TestCheckUserPasswordNotMatch(t *testing.T) {
	backend.InitPostgreSQLBackend()
	result, err := CheckUser("alice@alice", "bob")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != false {
		t.Errorf("Result should be false, but get true")
	}
}

func TestCheckUserEmailNotFound(t *testing.T) {
	backend.InitPostgreSQLBackend()
	result, err := CheckUser("bob@alice", "bob")
	if !errors.Is(err, util.ErrUserNotFound) {
		t.Errorf("Expect ErrUserNotFound, but get %v", err)
	}
	if result != false {
		t.Errorf("Result should be false, but get true")
	}
}

func TestAddUserGood(t *testing.T) {
	backend.InitPostgreSQLBackend()
	user := model.User{Email: "bob@bob", Username: "Bob", Password: "bob"}
	success, err := AddUser(&user)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if success == false {
		t.Errorf("Result should be true, but get false")
	}
}

func TestAddUserUserAlreadyExist(t *testing.T) {
	backend.InitPostgreSQLBackend()
	user := model.User{Email: "alice@alice", Username: "Alice2", Password: "alice2"}

	success, err := AddUser(&user)

	if err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	if success != false {
		t.Errorf("Result should be false, but get true")
	}
}

func TestGetUserIdGood(t *testing.T) {
	backend.InitPostgreSQLBackend()
	result, err := GetUserId("alice@alice")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 1 {
		t.Errorf("Result should be 1, but get %d", result)
	}
}

func TestGetUserIdEmailNotFound(t *testing.T) {
	backend.InitPostgreSQLBackend()
	result, err := GetUserId("bob@alice")
	if !errors.Is(err, util.ErrUserNotFound) {
		t.Errorf("Expect ErrUserNotFound, but get %v", err)
	}
	if result != 0 {
		t.Errorf("Result should be 0, but get %d", result)
	}
}
