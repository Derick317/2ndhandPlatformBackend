package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"secondHand/model"
	"secondHand/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddItemHandler(c *gin.Context) {
	// Process request
	email, _ := c.Get("email")
	fmt.Println(email)

	sellerId, err := service.GetUserId(email.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		Price:       float32(price),
		Tag:         model.TagType(tag),
		Description: c.PostForm("description"),
	}
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
}
