package main

import (
	"log"

	"github.com/LanangDepok/project-management/config"
	"github.com/LanangDepok/project-management/database/seed"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	log.Println("Running all seeders...")
	seed.SeedAdmin()
	log.Println("All seeders completed.")
}
