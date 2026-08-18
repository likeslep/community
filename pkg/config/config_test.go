package config

import (
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	type nested struct {
		Addr    string        `env:"TEST_ADDR" default:":8080"`
		Workers int           `env:"TEST_WORKERS" default:"4"`
		Verbose bool          `env:"TEST_VERBOSE" default:"false"`
		Timeout time.Duration `env:"TEST_TIMEOUT" default:"5s"`
	}
	t.Setenv("TEST_ADDR", ":9090")
	t.Setenv("TEST_WORKERS", "8")
	t.Setenv("TEST_VERBOSE", "true")
	t.Setenv("TEST_TIMEOUT", "10s")

	var c nested
	if err := Load(&c); err != nil {
		t.Fatalf("Load() err = %v", err)
	}
	if c.Addr != ":9090" || c.Workers != 8 || !c.Verbose || c.Timeout != 10*time.Second {
		t.Fatalf("Load() = %+v", c)
	}
}

func TestLoadDefault(t *testing.T) {
	type c struct {
		Addr string `env:"TEST_DEFAULT_ADDR_XYZ" default:":8080"`
	}
	var out c
	if err := Load(&out); err != nil {
		t.Fatalf("Load() err = %v", err)
	}
	if out.Addr != ":8080" {
		t.Fatalf("默认值未生效，got %q", out.Addr)
	}
}

func TestLoadRequiredMissing(t *testing.T) {
	type c struct {
		DSN string `env:"TEST_MISSING_DSN_XYZ" required:"true"`
	}
	var out c
	if err := Load(&out); err == nil {
		t.Fatal("Load() 期望缺失必填变量报错，实际 nil")
	}
}

func TestLoadNonPointer(t *testing.T) {
	var s string
	if err := Load(s); err == nil {
		t.Fatal("Load() 期望非指针报错，实际 nil")
	}
}
