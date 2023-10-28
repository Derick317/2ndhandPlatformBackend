package main

import (
	"fmt"

	"secondHand/backend"
	"secondHand/handler"
)

func main() {
	fmt.Println("started-service")
	backend.InitPostgreSQLBackend()
	backend.InitGCSBackend()
	handler.InitRouter()
	handler.RunRouter()
}
