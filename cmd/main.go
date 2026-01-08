package main

import (
	"log"
	"user_advt/internal/config"
)

func main() {
	config.MustLoad()

	log.Println("Config loaded")

}