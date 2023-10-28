package handler

import "github.com/gin-gonic/gin"

var router *gin.Engine

func InitRouter() {
	router = gin.Default()
	router.POST("/signup", signupHandler)
	router.POST("/signin", signinHandler)
	router.POST("/additem", authMiddleware(), AddItemHandler)
}

func RunRouter(addr ...string) {
	router.Run(addr...)
}
