package main

import (
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/util"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Driver struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	ProfilePicture string      `json:"profilePicture"`
	CarPlate       string      `json:"carPlate"`
	PackageSlug    string      `json:"packageSlug"`
	Location       *Coordinate `json:"location"`
	Geohash        string      `json:"geohash"`
}

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func handleDriversWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to upgrade to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Printf("userID is required")
		return
	}

	packageSlug := r.URL.Query().Get("packageSlug")
	if packageSlug == "" {
		log.Printf("packageSlug is required")
		return
	}

	// Send driver registration confirmation
	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: Driver{
			ID:             userID,
			Name:           "Tiago",
			ProfilePicture: util.GetRandomAvatar(1),
			CarPlate:       "ABC1234",
			PackageSlug:    packageSlug,
		},
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("failed to write message: %v", err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("failed to read message: %v", err)
			break
		}
		log.Printf("received message: %s", message)
	}
}

func handleRidersWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Printf("userID is required")
		return
	}

	// Send available drivers to rider
	drivers := []Driver{
		{
			ID:             "driver-1",
			Name:           "Tiago",
			ProfilePicture: util.GetRandomAvatar(1),
			CarPlate:       "ABC1234",
			PackageSlug:    "sedan",
			Location:       &Coordinate{Latitude: 37.7749, Longitude: -122.4194},
			Geohash:        "9q8yy",
		},
		{
			ID:             "driver-2",
			Name:           "Maria",
			ProfilePicture: util.GetRandomAvatar(2),
			CarPlate:       "XYZ5678",
			PackageSlug:    "suv",
			Location:       &Coordinate{Latitude: 37.7755, Longitude: -122.4180},
			Geohash:        "9q8yy",
		},
	}

	msg := contracts.WSMessage{
		Type: "driver.cmd.location",
		Data: drivers,
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("failed to write message: %v", err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("failed to read message: %v", err)
			break
		}
		log.Printf("received message: %s", message)
	}
}
