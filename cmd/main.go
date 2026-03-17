package main

import (
	"context"
	"net/http"
	"subscriptions-api/internal/config"
	"subscriptions-api/internal/database"
	"subscriptions-api/internal/handler"
	"subscriptions-api/internal/logger"
	"subscriptions-api/internal/repository"
	"subscriptions-api/internal/service"

	"github.com/gorilla/mux"

	_ "subscriptions-api/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Subscriptions API
// @version         1.0
// @description     CRUDL subscription service

// @BasePath  /subscriptions

func main() {
	ctx := context.Background()
	config, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("[MAIN] Error with config -> %w", err)
	}

	database, err := database.Connect(ctx, config)
	if err != nil {
		logger.Fatal("[MAIN] Error with database -> %w", err)
	}
	defer database.Close(ctx)

	repo := repository.NewSubscriptionRepository(database)
	service := service.NewSubscriptionRepository(repo)
	handler := handler.NewSubscriptionHandler(service)

	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.HandleFunc("/subscriptions/stats", handler.CollectStats).Methods("GET")
	router.HandleFunc("/subscriptions/{id}", handler.GetSubscription).Methods("GET")
	router.HandleFunc("/subscriptions/{id}", handler.UpdateSubscription).Methods("PATCH")
	router.HandleFunc("/subscriptions/{id}", handler.DeleteSubscription).Methods("DELETE")
	router.HandleFunc("/subscriptions", handler.ListSubscriptions).Methods("GET")
	router.HandleFunc("/subscriptions", handler.CreateSubscription).Methods("POST")

	http.ListenAndServe(":"+config.AppPort, router)
}
