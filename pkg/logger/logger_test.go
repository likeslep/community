package logger

import (
	"testing"

	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{"json-info", Config{Level: "info", JSON: true}, false},
		{"console-debug", Config{Level: "debug", JSON: false}, false},
		{"invalid-level", Config{Level: "bogus"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, err := New(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("New() err = %v, wantErr = %v", err, tt.wantErr)
			}
			if err == nil {
				l.Info("hello", zap.String("k", "v"))
			}
		})
	}
}
