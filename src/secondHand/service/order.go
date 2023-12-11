package service

import (
	"secondHand/backend"
	"secondHand/model"
	"secondHand/util"
	"time"

	"gorm.io/gorm"
)

// QueryBuyerOrders returns the IDs of Item in the order of buyer whose ID is BUYER_ID
func QueryBuyerOrders(buyerId uint64, tx *gorm.DB) (map[uint64]int64, error) {
	var orders []model.Order
	if err := backend.ReadFromDBByKey(&orders, "buyer_id", buyerId, false, tx); err != nil {
		return nil, err
	}
	items := make(map[uint64]int64, 0)
	for _, order := range orders {
		if order.ExpTime > time.Now().Unix() {
			items[order.ItemId] = order.ExpTime - time.Now().Unix()
		} else {
			CancelOrder(order)
		}
	}
	return items, nil
}

// Cancel Order removes order from orders table
// and resets the status of the corresponding item to be available
func CancelOrder(order model.Order) error {
	tx := backend.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	_, ok, err := TestAndSetItemStatus(order.ItemId, model.OnOrder, model.Available, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if !ok {
		tx.Rollback()
		return util.ErrOrderNoExists
	}
	if err := backend.DeleteFromDBByPrimaryKey(&model.Order{}, order.ID, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
