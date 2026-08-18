package service

import (
	"testing"

	"github.com/likeslep/community/internal/moderation/model"
)

func TestDetectorCheck(t *testing.T) {
	d := NewDetector([]model.SensitiveWord{
		{Word: "暴力", Level: model.LevelBlock},
		{Word: "广告", Level: model.LevelReview},
		{Word: "spam", Level: model.LevelBlock},
	})

	tests := []struct {
		name    string
		text    string
		wantLen int
	}{
		{"无命中", "正常的技术内容", 0},
		{"命中一个", "这里有暴力内容", 1},
		{"命中多个", "暴力加广告", 2},
		{"英文大小写不敏感", "这是 SPAM 内容", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(d.Check(tt.text)); got != tt.wantLen {
				t.Fatalf("Check(%q) 命中数 = %d, want %d", tt.text, got, tt.wantLen)
			}
		})
	}
}

func TestClassify(t *testing.T) {
	tests := []struct {
		name string
		hits []Hit
		want string
	}{
		{"无命中通过", nil, "pass"},
		{"仅 review", []Hit{{Word: "广告", Level: model.LevelReview}}, "review"},
		{"有 block 拦截", []Hit{{Word: "广告", Level: model.LevelReview}, {Word: "暴力", Level: model.LevelBlock}}, "block"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Classify(tt.hits); got != tt.want {
				t.Fatalf("Classify() = %q, want %q", got, tt.want)
			}
		})
	}
}
