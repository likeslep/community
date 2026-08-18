// Package sanitize 提供 Markdown 净化（防 XSS，plan.md §42）。
package sanitize

import (
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
)

// SanitizeMarkdown 渲染 Markdown 为安全 HTML，并剥离脚本等危险标签。
// 未安全化的原始 HTML 会被移除，防止存储型 XSS。
func SanitizeMarkdown(md string) string {
	var buf strings.Builder
	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		return ""
	}
	policy := bluemonday.UGCPolicy()
	return policy.Sanitize(buf.String())
}
