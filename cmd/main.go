package main

import (
	"context"
	"fmt"
	"log"
	"subscriptions-api/internal/config"
	"subscriptions-api/internal/database"
)

func main() {
	ctx := context.Background()
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	database, err := database.Connect(ctx, config)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(database)
}
