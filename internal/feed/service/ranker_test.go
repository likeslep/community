package service

import (
	"testing"
	"time"
)

func TestRuleRankerScore(t *testing.T) {
	now := time.Now()
	ranker := RuleRanker{}

	fresh := FeedItem{ID: "a", PublishedAt: now, LikeCount: 0, ViewCount: 0, Tags: nil, UserTags: nil}
	old := FeedItem{ID: "b", PublishedAt: now.Add(-10 * time.Hour), LikeCount: 0, ViewCount: 0}
	if ranker.Score(fresh, now) <= ranker.Score(old, now) {
		t.Fatal("新鲜内容得分应高于旧内容")
	}

	tagged := FeedItem{ID: "c", PublishedAt: now, Tags: []string{"go"}, UserTags: []string{"go"}}
	untagged := FeedItem{ID: "d", PublishedAt: now}
	if ranker.Score(tagged, now) <= ranker.Score(untagged, now) {
		t.Fatal("命中标签的内容得分应更高")
	}
}

func TestFreshnessScore(t *testing.T) {
	now := time.Now()
	if freshnessScore(now, now) != 100 {
		t.Fatalf("刚发布的新鲜度应为 100，got %f", freshnessScore(now, now))
	}
	if freshnessScore(now.Add(-time.Hour), now) <= 0 {
		t.Fatal("一小时前的新鲜度应为正")
	}
}

func TestTagMatchScore(t *testing.T) {
	score := tagMatchScore([]string{"go", "kafka"}, []string{"go", "rust"})
	if score != 10 {
		t.Fatalf("命中 1 个标签应得 10 分，got %f", score)
	}
}
