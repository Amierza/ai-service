package handler

import (
	"net/http"

	"github.com/Amierza/ai-service/dto"
	"github.com/Amierza/ai-service/response"
	"github.com/Amierza/ai-service/service"
	"github.com/gin-gonic/gin"
)

type (
	ISummaryYHandler interface {
		GenerateSummary(ctx *gin.Context)
	}

	summaryHandler struct {
		summaryService service.ISummaryService
	}
)

func NewSummaryHandler(summaryService service.ISummaryService) *summaryHandler {
	return &summaryHandler{
		summaryService: summaryService,
	}
}

func (sh *summaryHandler) GenerateSummary(ctx *gin.Context) {
	err := sh.summaryService.GenerateSummary(ctx)
	if err != nil {
		status := mapErrorToStatus(err)
		res := response.BuildResponseFailed(dto.FAILED_GENERATE_SUMMARY_WITH_GPT_LLM, err.Error(), nil)
		ctx.AbortWithStatusJSON(status, res)
		return
	}

	res := response.BuildResponseSuccess(dto.SUCCESS_GENERATE_SUMMARY_WITH_GPT_LLM, nil)
	ctx.JSON(http.StatusOK, res)
}
