package handler

import (
	"errors"
	"fmt"
	"net/http"
	"secondHand/model"
	"secondHand/service"
	"secondHand/util"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func searchHandler(c *gin.Context) {
	tag, err := strconv.ParseUint(c.Query("tag"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			gin.H{"status": fmt.Sprintf("invalid tag: \"%s\"", c.Query("tag"))})
		return
	}
	var items []model.Item
	keywords := strings.Fields(c.Query("keywords"))
	if items, err = service.Search(keywords, tag); errors.Is(err, util.ErrItemNotFound) {
		c.JSON(http.StatusOK, make([]string, 0))
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var itemIds []uint64
	for _, item := range items {
		itemIds = append(itemIds, item.ID)
	}
	c.JSON(http.StatusOK, itemIds)
}
