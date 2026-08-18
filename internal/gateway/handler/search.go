package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	searchv1 "github.com/likeslep/community/api/gen/search/v1"
)

// Search 搜索。
func (h *Handler) Search(c *gin.Context) {
	keyword := c.Query("keyword")
	typeFilter := c.Query("type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	resp, err := h.search.Search(c.Request.Context(), &searchv1.SearchRequest{
		Keyword: keyword, Type: typeFilter, Page: int32(page), PageSize: int32(pageSize),
	})
	if err != nil {
		writeErr(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": resp})
}
