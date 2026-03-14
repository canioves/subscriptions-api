package main

import (
	"context"
	"log"
	"net/http"
	"subscriptions-api/internal/config"
	"subscriptions-api/internal/database"
	"subscriptions-api/internal/handler"
	"subscriptions-api/internal/repository"
	"subscriptions-api/internal/service"

	"github.com/gorilla/mux"
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
	defer database.Close(ctx)

	repo := repository.NewSubscriptionRepository(database)
	service := service.NewSubscriptionRepository(repo)
	handler := handler.NewSubscriptionHandler(service)

	router := mux.NewRouter()
	router.HandleFunc("/subscriptions", handler.CreateSubscription).Methods("POST")
	router.HandleFunc("/subscriptions", handler.ListSubscriptions).Methods("GET")
	router.HandleFunc("/subscriptions/{id}", handler.GetSubscription).Methods("GET")

	http.ListenAndServe(":"+config.AppPort, router)
}
