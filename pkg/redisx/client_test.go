package redisx

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
)

func TestPing(t *testing.T) {
	mr := miniredis.RunT(t)
	c := New(Config{Addr: mr.Addr()})
	defer c.Close()

	if err := c.Ping(context.Background()); err != nil {
		t.Fatalf("Ping() err = %v", err)
	}
}

func TestSetGet(t *testing.T) {
	mr := miniredis.RunT(t)
	c := New(Config{Addr: mr.Addr()})
	defer c.Close()

	ctx := context.Background()
	if err := c.Redis().Set(ctx, "k", "v", 0).Err(); err != nil {
		t.Fatalf("Set() err = %v", err)
	}
	got, err := c.Redis().Get(ctx, "k").Result()
	if err != nil || got != "v" {
		t.Fatalf("Get() = %q, %v", got, err)
	}
}
