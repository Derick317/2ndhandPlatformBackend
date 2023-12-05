package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"secondHand/constants"
	"secondHand/model"
	"secondHand/service"
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

	// Handle request
	err = service.AddItem(&item, files)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("server cannot add item: %s", err.Error())})
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
	if err = service.QueryItem(&item, itemId); err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("unable to read item: %s", err.Error())})
		return
	}
	item.ImageUrls.Remove(constants.NEXTKEY_KEY)

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
