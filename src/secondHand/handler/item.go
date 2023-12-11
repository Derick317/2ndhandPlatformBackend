package handler

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"secondHand/backend"
	"secondHand/model"
	"secondHand/service"
	"secondHand/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

func addItemHandler(c *gin.Context) {
	// Process request
	sellerId, err := getUserIdFromGinContent(c)
	if err != nil {
		return
	}

	price, err := strconv.ParseFloat(c.PostForm("price"), 32)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("invalid price: \"%s\"", c.PostForm("price"))})
		return
	}

	tag, err := strconv.Atoi(c.PostForm("tag"))
	if err != nil || model.TagType(tag) >= model.TagCounter {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("invalid tag: \"%s\"", c.PostForm("tag"))})
		return
	}

	item := model.Item{
		SellerId:    sellerId,
		Title:       c.PostForm("title"),
		Price:       float32(price),
		Tag:         model.TagType(tag),
		Description: c.PostForm("description"),
	}
	// fmt.Println(item)
	fhs := c.Request.MultipartForm.File["image"]
	var files = []multipart.File{}
	for _, header := range fhs {
		file, err := header.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				gin.H{"error": fmt.Sprintf("cannot open file %s: %s", header.Filename, err.Error())})
			return
		}
		files = append(files, file)
	}

	tx := backend.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Handle request
	if err = service.AddItem(&item, files, tx); err != nil && !errors.Is(err, util.ErrGCS) {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("server cannot add item: %s", err.Error())})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func deleteItemHandler(c *gin.Context) {
	sellerId, itemId, err := getUserIdAndItemId(c)
	if err != nil {
		return
	}
	var item model.Item
	if service.QueryItem(&item, itemId, nil); errors.Is(err, util.ErrItemNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Item %d does not exist.", itemId)})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if item.SellerId != sellerId {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	tx := backend.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, ok, err := service.TestAndSetItemStatus(itemId, model.Available, model.Deleted, tx)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		status, ok, err := service.TestAndSetItemStatus(itemId, model.Sold, model.Deleted, tx)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !ok {
			tx.Rollback()
			c.JSON(http.StatusBadRequest,
				gin.H{"status": fmt.Sprintf("cannot delete item whose status is %d", status)})
			return
		}
	}

	if err = service.DeleteItem(&item, tx); err != nil && !errors.Is(err, util.ErrGCS) {
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

func queryItemHandler(c *gin.Context) {
	itemId, err := strconv.ParseUint(c.Query("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": fmt.Sprintf("invalid item ID: \"%s\"", c.Query("item_id"))})
		return
	}

	var item model.Item
	if err = service.QueryItem(&item, itemId, nil); errors.Is(err, util.ErrItemNotFound) {
		c.JSON(http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("Item %d does not exist.", itemId)})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("unable to read item: %s", err.Error())})
		return
	}

	c.PureJSON(http.StatusOK, gin.H{
		"id":          item.ID,
		"seller_id":   item.SellerId,
		"title":       item.Title,
		"price":       item.Price,
		"tag":         item.Tag,
		"description": item.Description,
		"status":      item.Status,
		"image_urls":  item.ImageUrls,
	})
}
