// Package service 是 file-service 的业务逻辑层。
package service

import (
	"context"
	"io"
	"strings"

	"github.com/google/uuid"

	"github.com/likeslep/community/internal/file/model"
	"github.com/likeslep/community/internal/file/storage"
)

// 允许上传的 MIME 类型前缀。
var allowedTypes = []string{"image/", "text/markdown", "text/plain", "application/pdf"}

// Config 是 file-service 的业务配置。
type Config struct {
	MaxSize int64
	Storage storage.Storage
}

// FileService 是文件业务逻辑层。
type FileService struct {
	repo Repository
	cfg  Config
}

// NewFileService 构造。
func NewFileService(repo Repository, cfg Config) *FileService {
	return &FileService{repo: repo, cfg: cfg}
}

// Upload 上传文件：校验 → 存储 → 记录元数据。
func (s *FileService) Upload(ctx context.Context, userID uint64, name, contentType string, size int64, r io.Reader) (*model.File, error) {
	if size > s.cfg.MaxSize {
		return nil, ErrTooLarge
	}
	if !isAllowed(contentType) {
		return nil, ErrInvalidType
	}
	key := uuid.NewString()
	if _, err := s.cfg.Storage.Put(ctx, key, r); err != nil {
		return nil, err
	}
	f := &model.File{UserID: userID, Name: name, Path: key, Type: contentType, Size: size}
	if err := s.repo.Create(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

// Download 下载文件：查询元数据 → 读取存储。
func (s *FileService) Download(ctx context.Context, id uint64) (*model.File, io.ReadCloser, error) {
	f, err := s.repo.Find(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	r, err := s.cfg.Storage.Get(ctx, f.Path)
	if err != nil {
		return nil, nil, err
	}
	return f, r, nil
}

func isAllowed(contentType string) bool {
	for _, t := range allowedTypes {
		if strings.HasPrefix(contentType, t) {
			return true
		}
	}
	return false
}
