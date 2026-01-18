package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/env"
)

var (
	tripServiceURL = env.GetString("TRIP_SERVICE_URL", "http://trip-service:8083")
)

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Second * 9)

	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validation
	if reqBody.UserID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	// Forward request to trip-service
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		http.Error(w, "failed to marshal request", http.StatusInternalServerError)
		return
	}

	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	defer tripService.Close()

	// tripService.Client.PreviewTrip()

	tripResp, err := http.Post(
		fmt.Sprintf("%s/trip/preview", tripServiceURL),
		"application/json",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to call trip service: %v", err), http.StatusInternalServerError)
		return
	}
	defer tripResp.Body.Close()

	// Read response from trip-service
	body, err := io.ReadAll(tripResp.Body)
	if err != nil {
		http.Error(w, "failed to read trip service response", http.StatusInternalServerError)
		return
	}

	// Forward the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(tripResp.StatusCode)
	w.Write(body)
}
