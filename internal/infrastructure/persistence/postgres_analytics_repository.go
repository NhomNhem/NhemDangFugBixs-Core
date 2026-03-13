package persistence

import (
	"context"
	"encoding/json"

	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/models"
	"github.com/NhomNhem/HollowWilds-Backend/internal/domain/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresAnalyticsRepository struct {
	db *pgxpool.Pool
}

// NewPostgresAnalyticsRepository creates a new PostgreSQL analytics repository
func NewPostgresAnalyticsRepository(db *pgxpool.Pool) repository.AnalyticsRepository {
	return &postgresAnalyticsRepository{db: db}
}

func (r *postgresAnalyticsRepository) RecordEvents(ctx context.Context, events []models.AnalyticsEvent) (int, int, error) {
	if r.db == nil {
		return len(events), 0, nil
	}
	accepted := 0
	rejected := 0

	for _, event := range events {
		payloadJSON, _ := json.Marshal(event.Payload)

		_, err := r.db.Exec(ctx, `
			INSERT INTO analytics_events (user_id, session_id, event_type, event_properties, created_at)
			VALUES ($1, $2, $3, $4, NOW())
		`, event.UserID, event.SessionID, event.EventName, payloadJSON)

		if err != nil {
			rejected++
		} else {
			accepted++
		}
	}

	return accepted, rejected, nil
}
