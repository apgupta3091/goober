package service

import (
	"context"
	"ride-sharing/services/trip-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type TripModel struct {
// 	ID       primitive.ObjectID
// 	UserID   string
// 	Status   string
// 	RideFare RideFareModel
// }

// type TripRepository interface {
// 	CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
// }

// type TripService interface {
// 	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
// }

type Service struct {
	repo domain.TripRepository
}

func NewService(repo domain.TripRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateTrip(ctx context.Context, fare *domain.RideFareModel) (*domain.TripModel, error) {
	t := &domain.TripModel{
		ID:       primitive.NewObjectID(),
		UserID:   fare.UserID,
		Status:   "pending",
		RideFare: fare,
	}
	return s.repo.CreateTrip(ctx, t)
}
