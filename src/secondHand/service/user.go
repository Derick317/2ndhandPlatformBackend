package service

import "secondHand/model"

// AddUser adds a user whose infomation is represented as a instance of model.User
// It reports whether the user is added successfully
func AddUser(user *model.User) (bool, error) {
	return false, nil
}

func CheckUser(email, password string) (bool, error) {
	return false, nil
}
