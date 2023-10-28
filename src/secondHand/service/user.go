package service

import (
	"errors"
	"secondHand/backend"
	"secondHand/model"
	"secondHand/util"

	"gorm.io/gorm"
)

// AddUser adds a user whose infomation is represented as a instance of model.User
// It reports whether the user is added successfully.
func AddUser(user *model.User) (bool, error) {
	// Checker whether the email already exists
	var user_temp model.User
	if backend.ReadFromDBByKey(&user_temp, "email", user.Email, true) == nil {
		return false, nil
	}
	if err := backend.CreateRecord(user); err != nil {
		return false, err
	}
	return true, nil
}

// Checkuser checks whether PASSWORD matches the true password in the database.
// Possible errors:
//   - ErrUserNotFound: the email has not been registed yet
func CheckUser(email, password string) (bool, error) {
	var user model.User
	if err := backend.ReadFromDBByKey(&user, "email", email, true); err == nil {
		return password == user.Password, nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, util.ErrUserNotFound
	} else {
		return false, err
	}
}

func GetUserId(email string) (uint, error) {
	var user model.User
	if err := backend.ReadFromDBByKey(&user, "email", email, true); err == nil {
		return user.ID, nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, util.ErrUserNotFound
	} else {
		return 0, err
	}
}
