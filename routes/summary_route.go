package routes

import (
	"github.com/Amierza/ai-service/handler"
	"github.com/Amierza/ai-service/jwt"
	"github.com/Amierza/ai-service/middleware"
	"github.com/gin-gonic/gin"
)

func Summary(route *gin.Engine, summaryHandler handler.ISummaryYHandler, jwt jwt.IJWT) {
	routes := route.Group("/api/v1").Use(middleware.Authentication(jwt))
	{
		routes.GET("/generate-summary", summaryHandler.GenerateSummary)
	}
}
