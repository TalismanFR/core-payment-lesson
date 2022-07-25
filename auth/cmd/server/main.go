package main

import (
	"auth/internal/cmd"
	"auth/internal/config"
	"log"
)

func main() {

	cfg := config.Init("configs/config.yaml")

	if err := cmd.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
