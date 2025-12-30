package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RideFareModel struct {
	ID                primitive.ObjectID
	UserID            string
	PackageSlug       string //ex van, luxury, etc.
	TotalPriceInCents float64
	ExpiresAt         time.Time
}
