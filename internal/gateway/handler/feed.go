package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	feedv1 "github.com/likeslep/community/api/gen/feed/v1"
	gwmw "github.com/likeslep/community/internal/gateway/middleware"
	"github.com/likeslep/community/pkg/grpcx"
)

// GetFeed 查询信息流。
func (h *Handler) GetFeed(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	ctx := grpcx.WithUserID(c.Request.Context(), c.GetString(gwmw.CtxUserID))
	resp, err := h.feed.GetFeed(ctx, &feedv1.GetFeedRequest{Limit: int32(limit)})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}
