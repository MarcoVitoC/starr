package main

import "github.com/MarcoVitoC/starr/internal"

func main() {
	server := internal.NewServer("localhost:8080")
	server.Run()
}