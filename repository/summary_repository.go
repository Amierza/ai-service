package repository

import (
	"gorm.io/gorm"
)

type (
	ISummaryRepository interface {
		// CREATE / POST

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
