package repository

import (
	"context"

	"github.com/Amierza/ai-service/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	ISummaryRepository interface {
		// CREATE / POST
		SaveSummary(ctx context.Context, tx *gorm.DB, sessionID, summary string) error

		// READ / GET

		// UPDATE / PATCH

		// DELETE / DELETE
	}

	summaryRepository struct {
		db *gorm.DB
	}
)

func NewSummaryRepository(db *gorm.DB) *summaryRepository {
	return &summaryRepository{
		db: db,
	}
}

func (sr *summaryRepository) SaveSummary(ctx context.Context, tx *gorm.DB, sessionID, summary string) error {
	if tx == nil {
		tx = sr.db
	}

	// Konversi sessionID dari string ke UUID
	sid, err := uuid.Parse(sessionID)
	if err != nil {
		return err
	}

	note := entity.Note{
		ID:        uuid.New(),
		Content:   summary,
		SessionID: sid,
	}

	// Simpan ke database
	if err := tx.WithContext(ctx).Create(&note).Error; err != nil {
		return err
	}

	return nil
}
