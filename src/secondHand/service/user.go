package service

import (
	"errors"
	"fmt"
	"secondHand/backend"
	"secondHand/constants"
	"secondHand/model"
	"secondHand/util"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// AddUser adds a user whose infomation is represented as a instance of model.User
// It reports whether the user is added successfully.
func AddUser(user *model.User, tx *gorm.DB) (bool, error) {
	// Checker whether the email already exists
	var user_temp model.User
	if backend.ReadFromDBByKey(&user_temp, "email", user.Email, true, tx) == nil {
		return false, nil
	}
	if err := backend.CreateRecord(user, tx); err != nil {
		return false, err
	}
	return true, nil
}

// Checkuser checks whether PASSWORD matches the true password in the database.
// USER is the entity to store the user's information if successful.
// Possible errors:
//   - ErrUserNotFound: the email has not been registed yet
func CheckUser(user *model.User, email, password string, tx *gorm.DB) (bool, error) {
	if err := backend.ReadFromDBByKey(user, "email", email, true, tx); err == nil {
		return password == user.Password, nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, util.ErrUserNotFound
	} else {
		return false, err
	}
}

func GetUserId(email string, tx *gorm.DB) (uint64, error) {
	var user model.User
	if err := backend.ReadFromDBByKey(&user, "email", email, true, tx); err == nil {
		return user.ID, nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, util.ErrUserNotFound
	} else {
		return 0, err
	}
}

// UserAddOrder add USER_ID-th user and ITEM_ID-th item to the table.
func UserAddOrder(userId uint64, itemId uint64, tx *gorm.DB) error {
	order := model.Order{
		BuyerId: userId,
		ItemId:  itemId,
		ExpTime: time.Now().Add(constants.ORDER_EXPIRE_TIME).Unix(),
	}

	return backend.CreateRecord(&order, tx)
}

// UserRemoveOrder removes user's ITEM_ID-th item from the order table.
// It reports whether it is successfully removed.
// Failure happends if item does not exist.
func UserRemoveOrder(userId uint64, itemId uint64, tx *gorm.DB) (bool, error) {
	num, err := backend.DeleteFromDBByKeys(&model.Order{}, []string{"buyer_id", "item_id"},
		[]string{strconv.FormatUint(userId, 10), strconv.FormatUint(itemId, 10)}, tx)
	if err != nil {
		return false, err
	}
	if num == 0 {
		return false, nil
	}
	if num > 1 {
		return false, fmt.Errorf("%w: delete %d orders but expect only one", util.ErrUnexpected, num)
	}
	return true, nil
}

// QuerySellerList returns the IDs of items in the list of seller whose ID is SELLER_ID
func QuerySellerList(sellerId uint64, tx *gorm.DB) ([]uint64, error) {
	var items []model.Item
	if err := backend.ReadFromDBByKey(&items, "seller_id", sellerId, false, tx); err != nil {
		return nil, err
	}
	list := []uint64{}
	for _, item := range items {
		list = append(list, item.ID)
	}
	return list, nil
}
