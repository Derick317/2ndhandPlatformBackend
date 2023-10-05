package main

import (
    "fmt"

    "secondHand/backend"
)

func main() {
    fmt.Println("started-service")
    backend.InitPostgreSQLBackend()
}