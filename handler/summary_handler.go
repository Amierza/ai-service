package handler

import (
	pb "github.com/Amierza/ai-service/proto"
	"github.com/Amierza/ai-service/service"
	"github.com/gin-gonic/gin"
)

type (
	ISummaryYHandler interface {
		GenerateSummary(ctx *gin.Context, req *pb.SummaryRequest) (*pb.SummaryResponse, error)
	}

	summaryHandler struct {
		pb.UnimplementedSummaryServiceServer
		summaryService service.ISummaryService
	}
)

func NewSummaryHandler(summaryService service.ISummaryService) *summaryHandler {
	return &summaryHandler{
		summaryService: summaryService,
	}
}

func (sh *summaryHandler) GenerateSummary(ctx *gin.Context, req *pb.SummaryRequest) (*pb.SummaryResponse, error) {
	result, err := sh.summaryService.GenerateSummary(ctx, req)
	if err != nil {
		return nil, err
	}

	return &pb.SummaryResponse{
		SessionId: req.Task.SessionId,
		Summary:   result,
		Status:    "success",
	}, nil
}
