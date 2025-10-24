package service

import (
	"context"

	pb "github.com/Amierza/ai-service/proto"
	"go.uber.org/zap"
)

type GRPCSummaryServer struct {
	pb.UnimplementedSummaryServiceServer
	summaryService ISummaryService
	logger         *zap.Logger
}

func NewGRPCSummaryServer(summaryService ISummaryService, logger *zap.Logger) *GRPCSummaryServer {
	return &GRPCSummaryServer{
		summaryService: summaryService,
		logger:         logger,
	}
}

func (s *GRPCSummaryServer) GenerateSummary(ctx context.Context, req *pb.SummaryRequest) (*pb.SummaryResponse, error) {
	s.logger.Info("🧠 Received gRPC GenerateSummary request", zap.String("session_id", req.Task.SessionId))

	summary, err := s.summaryService.GenerateSummary(ctx, req)
	if err != nil {
		s.logger.Error("failed to generate summary", zap.Error(err))
		return nil, err
	}

	return &pb.SummaryResponse{
		Summary: summary,
	}, nil
}
