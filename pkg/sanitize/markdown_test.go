package sanitize

import (
	"strings"
	"testing"
)

func TestSanitizeMarkdown(t *testing.T) {
	// 危险脚本被移除
	out := SanitizeMarkdown(`<script>alert('xss')</script>`)
	if strings.Contains(out, "<script") {
		t.Fatal("script 标签应被移除")
	}
	// Markdown 渲染为 HTML
	html := SanitizeMarkdown("**bold**")
	if !strings.Contains(html, "<strong>") {
		t.Fatalf("Markdown 应渲染为 <strong>，got %q", html)
	}
}

func TestSanitizePlainText(t *testing.T) {
	out := SanitizeMarkdown("普通文本 <img src=x onerror=alert(1)>")
	if strings.Contains(out, "onerror") {
		t.Fatal("危险属性应被移除")
	}
}
