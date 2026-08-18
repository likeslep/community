// Package service 是 moderation-service 的业务逻辑层。
package service

import (
	"strings"

	"github.com/likeslep/community/internal/moderation/model"
)

// Detector 是敏感词检测器，基于大小写不敏感的子串匹配。
// 注：当前为线性扫描实现，词库规模增大后可替换为 Aho-Corasick 自动机。
type Detector struct {
	words []model.SensitiveWord
}

// NewDetector 构造检测器。
func NewDetector(words []model.SensitiveWord) *Detector {
	return &Detector{words: words}
}

// Hit 是一次敏感词命中。
type Hit struct {
	Word  string
	Level string
}

// Check 返回文本命中的所有敏感词。
func (d *Detector) Check(text string) []Hit {
	lower := strings.ToLower(text)
	var hits []Hit
	for _, w := range d.words {
		if strings.Contains(lower, strings.ToLower(w.Word)) {
			hits = append(hits, Hit{Word: w.Word, Level: w.Level})
		}
	}
	return hits
}

// Classify 根据命中结果返回处理建议：pass / review / block（plan.md §25）。
func Classify(hits []Hit) string {
	action := "pass"
	for _, h := range hits {
		if h.Level == model.LevelBlock {
			return "block"
		}
		action = "review"
	}
	return action
}
