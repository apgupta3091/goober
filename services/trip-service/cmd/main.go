package main

import (
	"log"
	"net/http"
	triphttp "ride-sharing/services/trip-service/internal/infrastructure/http"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/env"
)

var (
	httpAddr = env.GetString("HTTP_ADDR", ":8083")
)

func main() {
	log.Println("Starting Trip Service")

	inMemRepo := repository.NewInMemRepository()
	svc := service.NewService(inMemRepo)

	handler := &triphttp.HttpHandler{
		Service: svc,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /trip/preview", handler.HandleTripPreview)

	server := &http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	log.Printf("Trip Service listening on %s", httpAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
