package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func InitRouter() {
	router = gin.Default()

	router.Use(cors())
	router.POST("/signup", signupHandler)
	router.POST("/signin", signinHandler)
	router.GET("/search", searchHandler)
	router.GET("/qitem", queryItemHandler)
	router.GET("/qlist", authMiddleware(), querySellerListHandler)
	router.POST("/additem", authMiddleware(), addItemHandler)
	router.DELETE("/ditem", authMiddleware(), deleteItemHandler)
	router.POST("/addorder", authMiddleware(), addOrderHandler)
	router.POST("/cancelorder", authMiddleware(), cancelOrderHander)
	router.POST("/checkout", authMiddleware(), payForOrderHandler)
	router.GET("/qorder", authMiddleware(), queryOrderHandler)
}

func RunRouter(addr ...string) {
	router.Run(addr...)
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Methods", "POST, GET, DELETE")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
