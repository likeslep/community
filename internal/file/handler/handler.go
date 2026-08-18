// Package handler 是 file-service 的 HTTP 处理器（文件二进制走 HTTP，plan.md §26）。
package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/likeslep/community/internal/file/service"
)

// Handler 是 file-service 的 HTTP 处理器。
type Handler struct {
	svc *service.FileService
}

// New 构造。
func New(svc *service.FileService) *Handler { return &Handler{svc: svc} }

// RegisterRoutes 注册路由。
func (h *Handler) RegisterRoutes(engine *gin.Engine) {
	engine.POST("/files", h.Upload)
	engine.GET("/files/:id", h.Download)
}

// Upload 处理文件上传。
func (h *Handler) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少文件字段"})
		return
	}
	defer file.Close()

	f, err := h.svc.Upload(c.Request.Context(), 0, header.Filename, header.Header.Get("Content-Type"), header.Size, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"file_id": f.ID}})
}

// Download 处理文件下载。
func (h *Handler) Download(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "非法的文件 ID"})
		return
	}
	f, r, err := h.svc.Download(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "文件不存在"})
		return
	}
	defer r.Close()

	c.Header("Content-Type", f.Type)
	c.Header("Content-Disposition", "attachment; filename="+f.Name)
	c.Header("Content-Length", strconv.FormatInt(f.Size, 10))
	if _, err := io.Copy(c.Writer, r); err != nil {
		c.Error(err)
	}
}
