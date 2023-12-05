package service

import (
	"errors"
	"mime/multipart"
	"secondHand/backend"
	"secondHand/model"
	"secondHand/util"
	"testing"

	"gorm.io/gorm"
)

func TestCheckUserMatch(t *testing.T) {
	backend.InitPostgreSQLBackend()
	var user model.User
	result, err := CheckUser(&user, "alice@alice", "alice")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result == false {
		t.Errorf("Result should be true, but get false")
	}
	if user.ID != 1 {
		t.Errorf("Expected ID = 1, but get %d", user.ID)
	}
}

func TestCheckUserPasswordNotMatch(t *testing.T) {
	backend.InitPostgreSQLBackend()
	var user model.User
	result, err := CheckUser(&user, "alice@alice", "bob")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != false {
		t.Errorf("Result should be false, but get true")
	}
}

func TestCheckUserEmailNotFound(t *testing.T) {
	backend.InitPostgreSQLBackend()
	var user model.User
	result, err := CheckUser(&user, "bob@alice", "bob")
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

func TestUserAddOrder(t *testing.T) {
	backend.InitPostgreSQLBackend()
	var user model.User
	if err := backend.ReadFromDBByKey(&user, "username", "Alice", true); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	if err := UserAddOrder(user.ID, 2); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	if err := backend.ReadFromDBByKey(&user, "username", "Alice", true); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestUserAddTwoOrders(t *testing.T) {
	backend.InitPostgreSQLBackend()
	addTwoOrders(t)
}

func TestUserRemoveOrders(t *testing.T) {
	backend.InitPostgreSQLBackend()
	var user model.User
	addTwoOrders(t)
	if ok, err := UserRemoveOrder(user.ID, 2); err != nil {
		t.Errorf("Unexpect error: %v", err)
	} else if !ok {
		t.Errorf("Fail to remove order!")
	}
	var order model.Order
	err := backend.ReadFromDBByKey(&order, "item_id", "2", true)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Unexpect error: %v", err)
	}
	err = backend.ReadFromDBByKey(&order, "item_id", "10", true)
	if err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
}

func TestQuerySellerList(t *testing.T) {
	var item1 = model.Item{SellerId: 2, Price: 10.4, Tag: 1, Description: "Total football"}
	var item2 = model.Item{SellerId: 2, Price: 0.5, Tag: 0, Description: "Baseball"}
	backend.InitPostgreSQLBackend()
	if err := AddItem(&item1, []multipart.File{}); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	if err := AddItem(&item2, []multipart.File{}); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	if list, err := QuerySellerList(2); err != nil {
		t.Errorf("Unexpect error: %v", err)
	} else if len(list) != 2 || !(list[0] == 1 && list[1] == 2 || list[0] == 2 && list[1] == 1) {
		t.Errorf("Unexpect list: %v", list)
	}
}

// addTwoOrders adds the 2nd item and the 10th item to Alice's order
func addTwoOrders(t *testing.T) {
	var user model.User
	if err := backend.ReadFromDBByKey(&user, "username", "Alice", true); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	if err := UserAddOrder(user.ID, 2); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	if err := backend.ReadFromDBByKey(&user, "username", "Alice", true); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
	if err := UserAddOrder(user.ID, 10); err != nil {
		t.Errorf("Unexpect error: %v", err)
	}
}
