package service

import (
	"context"

	"github.com/Amierza/ai-service/jwt"
	"github.com/Amierza/ai-service/repository"
	"go.uber.org/zap"
)

type (
	ISummaryService interface {
		GenerateSummary(ctx context.Context) error
	}

	summaryService struct {
		summaryRepo repository.ISummaryRepository
		logger      *zap.Logger

		jwt jwt.IJWT
	}
)

func NewSummaryService(summaryRepo repository.ISummaryRepository, logger *zap.Logger, jwt jwt.IJWT) *summaryService {
	return &summaryService{
		summaryRepo: summaryRepo,
		logger:      logger,

		jwt: jwt,
	}
}

func (ss *summaryService) GenerateSummary(ctx context.Context) error { return nil }
