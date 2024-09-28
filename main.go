package main

import (
	"fmt"
	"log"
)

func main() {
	store, err := NewPostGresStore()
	if err != nil {
		log.Fatal(err)
	}
	server := NewAPIServer(":8080", store)
	server.Run()
	fmt.Println("Server running!!")
}
