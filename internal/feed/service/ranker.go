package service

import "time"

// FeedItem 是信息流条目。
type FeedItem struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	AuthorID    uint64    `json:"author_id"`
	PublishedAt time.Time `json:"published_at"`
	LikeCount   int64     `json:"like_count"`
	ViewCount   int64     `json:"view_count"`
	Tags        []string  `json:"tags"`
	UserTags    []string  `json:"user_tags"` // 用户关注的标签（用于 tag_match）
}

// Ranker 是排序引擎抽象接口，未来可替换为推荐系统（plan.md §23）。
type Ranker interface {
	Score(item FeedItem, now time.Time) float64
}

// RuleRanker 是规则排序实现：score = freshness + engagement + author + tag_match。
type RuleRanker struct{}

// Score 计算条目得分。
func (RuleRanker) Score(item FeedItem, now time.Time) float64 {
	return freshnessScore(item.PublishedAt, now) +
		engagementScore(item.LikeCount, item.ViewCount) +
		authorScore(item.AuthorID) +
		tagMatchScore(item.Tags, item.UserTags)
}

// freshnessScore 内容新鲜度，按时间衰减（越新分越高）。
func freshnessScore(publishedAt, now time.Time) float64 {
	age := max(now.Sub(publishedAt), 0)
	hours := age.Hours()
	if hours < 1 {
		return 100
	}
	return 100 / hours
}

// engagementScore 互动热度。
func engagementScore(like, view int64) float64 {
	return float64(like)*2 + float64(view)*0.1
}

// authorScore 作者权重（简化：固定值，后续可扩展为作者等级）。
func authorScore(_ uint64) float64 {
	return 1
}

// tagMatchScore 标签兴趣匹配（命中用户关注标签越多分越高）。
func tagMatchScore(tags, userTags []string) float64 {
	userSet := make(map[string]struct{}, len(userTags))
	for _, t := range userTags {
		userSet[t] = struct{}{}
	}
	var score float64
	for _, t := range tags {
		if _, ok := userSet[t]; ok {
			score += 10
		}
	}
	return score
}
