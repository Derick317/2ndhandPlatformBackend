package main

import (
	"fmt"

	"secondHand/backend"
	"secondHand/handler"
	"secondHand/service"
)

func main() {
	fmt.Println("started-2nd-hand-service")
	backend.InitPostgreSQLBackend()
	backend.InitGCSBackend()
	service.InitOrderCanceler()
	handler.InitRouter()
	handler.RunRouter()
}
