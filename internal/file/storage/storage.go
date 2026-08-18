// Package storage 提供文件存储抽象（plan.md §26：Storage 接口 + LocalStorage）。
package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

// Storage 是文件存储抽象接口，未来可替换为 OSS/S3/MinIO。
type Storage interface {
	Put(ctx context.Context, key string, r io.Reader) (int64, error)
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	Delete(ctx context.Context, key string) error
}

// LocalStorage 是本地文件系统实现。
type LocalStorage struct {
	baseDir string
}

// NewLocalStorage 构造本地存储。
func NewLocalStorage(baseDir string) *LocalStorage {
	return &LocalStorage{baseDir: baseDir}
}

// Put 写入文件，返回写入字节数。key 经过 Clean 防止路径穿越。
func (l *LocalStorage) Put(_ context.Context, key string, r io.Reader) (int64, error) {
	path := l.resolve(key)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return 0, err
	}
	f, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return io.Copy(f, r)
}

// Get 读取文件。
func (l *LocalStorage) Get(_ context.Context, key string) (io.ReadCloser, error) {
	return os.Open(l.resolve(key))
}

// Delete 删除文件。
func (l *LocalStorage) Delete(_ context.Context, key string) error {
	return os.Remove(l.resolve(key))
}

// resolve 将 key 解析为绝对路径，并限制在 baseDir 内。
func (l *LocalStorage) resolve(key string) string {
	return filepath.Join(l.baseDir, filepath.Clean("/"+key))
}
