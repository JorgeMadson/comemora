package repository

import (
	"celebrationhub/internal/core/domain"
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(dsn string) (*PostgresRepository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&domain.Event{}); err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) Save(ctx context.Context, event *domain.Event) error {
	return r.db.WithContext(ctx).Save(event).Error
}

func (r *PostgresRepository) FindByID(ctx context.Context, id uint) (*domain.Event, error) {
	var event domain.Event
	if err := r.db.WithContext(ctx).First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]domain.Event, error) {
	var events []domain.Event
	if err := r.db.WithContext(ctx).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Event{}, id).Error
}

func (r *PostgresRepository) FindByDate(ctx context.Context, day, month int) ([]domain.Event, error) {
	var events []domain.Event
	// Filter by Day and Month. Year is ignored for recurrent events.
	if err := r.db.WithContext(ctx).Where("day = ? AND month = ?", day, month).Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}
