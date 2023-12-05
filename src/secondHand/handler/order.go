package handler

import (
	"errors"
	"fmt"
	"net/http"
	"secondHand/backend"
	"secondHand/model"
	"secondHand/service"
	"secondHand/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

func addOrderHandler(c *gin.Context) {
	buyerId, itemId, err := getBuyerIdAndItemId(c)
	if err != nil {
		return
	}

	// Note the use of tx as the database handle once you are within a transaction
	tx := backend.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	status, ok, err := service.TestAndSetItemStatus(itemId, model.Available, model.OnOrder)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"status": fmt.Sprintf("cannot order: %d", status)})
		return
	}
	if err := service.UserAddOrder(buyerId, itemId); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func payForOrderHandler(c *gin.Context) {
	buyerId, itemId, err := getBuyerIdAndItemId(c)
	if err != nil {
		return
	}

	tx := backend.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if ok, err := service.UserRemoveOrder(buyerId, itemId); err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if !ok {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"status": "order does not exist"})
		return
	}
	status, ok, err := service.TestAndSetItemStatus(itemId, model.OnOrder, model.Sold)
	if err != nil || !ok {
		tx.Rollback()
		var errStr string
		if err != nil {
			errStr = err.Error()
		} else {
			errStr = fmt.Sprintf("cannot handle order: %d", status)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": errStr})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func queryOrderHandler(c *gin.Context) {
	buyerId, err := getUserIdFromGinContent(c)
	if err != nil {
		return
	}
	orders, err := service.QueryBuyerOrders(buyerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func cancelOrderHander(c *gin.Context) {
	buyerId, itemId, err := getBuyerIdAndItemId(c)
	if err != nil {
		return
	}
	var order model.Order
	if err := backend.ReadFromDBByKeys(&order, []string{"item_id", "buyer_id"},
		[]string{strconv.FormatUint(itemId, 10), strconv.FormatUint(buyerId, 10)}, true); err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "cannot read information of orders from database" + err.Error()})
		return
	}
	if err := service.CancelOrder(order); errors.Is(err, util.ErrOrderNoExists) {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": "order does not exist!"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// "item_id" should be sent in a postform
func getBuyerIdAndItemId(c *gin.Context) (uint64, uint64, error) {
	buyerId, err := getUserIdFromGinContent(c)
	if err != nil {
		return 0, 0, err
	}
	itemId, err := strconv.ParseUint(c.PostForm("item_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("invalid itemId: %s", c.PostForm("item_id"))})
		return 0, 0, err
	}
	return buyerId, itemId, nil
}
