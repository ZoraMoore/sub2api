package routes

import (
	"github.com/Wei-Shaw/sub2api/internal/handler"

	"github.com/gin-gonic/gin"
)

// RegisterMarketplaceRoutes 注册公开模型广场路由（无需认证）。
func RegisterMarketplaceRoutes(v1 *gin.RouterGroup, h *handler.Handlers) {
	marketplace := v1.Group("/marketplace")
	{
		marketplace.GET("/models", h.ModelMarketplace.ListPublic)
		marketplace.GET("/stats", h.ModelMarketplace.StatsPublic)
	}
}
