package storage

import (
	"context"
	"io"
	"strings"
	"testing"
)

func TestLocalStorage(t *testing.T) {
	ctx := context.Background()
	s := NewLocalStorage(t.TempDir())

	n, err := s.Put(ctx, "dir/file.txt", strings.NewReader("hello world"))
	if err != nil {
		t.Fatalf("Put() err = %v", err)
	}
	if n != 11 {
		t.Fatalf("写入字节数 = %d, want 11", n)
	}

	r, err := s.Get(ctx, "dir/file.txt")
	if err != nil {
		t.Fatalf("Get() err = %v", err)
	}
	data, err := io.ReadAll(r)
	r.Close()
	if err != nil || string(data) != "hello world" {
		t.Fatalf("读取内容 = %q, %v", data, err)
	}

	if err := s.Delete(ctx, "dir/file.txt"); err != nil {
		t.Fatalf("Delete() err = %v", err)
	}
	if _, err := s.Get(ctx, "dir/file.txt"); err == nil {
		t.Fatal("删除后应读取失败")
	}
}
