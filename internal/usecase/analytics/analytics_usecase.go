package analytics

import (
	"context"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/repository"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/usecase"
	"github.com/google/uuid"
)

type analyticsUsecase struct {
	repo repository.AnalyticsRepository
}

// NewAnalyticsUsecase creates a new analytics usecase
func NewAnalyticsUsecase(repo repository.AnalyticsRepository) usecase.AnalyticsUsecase {
	return &analyticsUsecase{repo: repo}
}

func (u *analyticsUsecase) TrackEvents(ctx context.Context, playerID *uuid.UUID, events []models.AnalyticsEvent) (int, int, error) {
	// Add playerID to events if present
	if playerID != nil {
		for i := range events {
			events[i].UserID = playerID
		}
	}

	return u.repo.RecordEvents(ctx, events)
}
