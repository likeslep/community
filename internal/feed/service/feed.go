package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/likeslep/community/pkg/redisx"
)

// FeedService 是信息流业务逻辑层，Redis Sorted Set 存储（plan.md §23）。
type FeedService struct {
	rdb    *redisx.Client
	ranker Ranker
}

// NewFeedService 构造。
func NewFeedService(rdb *redisx.Client, ranker Ranker) *FeedService {
	return &FeedService{rdb: rdb, ranker: ranker}
}

// AddToFeed 将条目写入用户 feed（按排序得分）。
func (s *FeedService) AddToFeed(ctx context.Context, userID uint64, item FeedItem) error {
	score := s.ranker.Score(item, time.Now())
	data, err := json.Marshal(item)
	if err != nil {
		return err
	}
	return s.rdb.Redis().ZAdd(ctx, feedKey(userID), redis.Z{Score: score, Member: string(data)}).Err()
}

// GetFeed 读取用户 feed（按得分降序）。
func (s *FeedService) GetFeed(ctx context.Context, userID uint64, limit int) ([]FeedItem, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	members, err := s.rdb.Redis().ZRangeArgs(ctx, redis.ZRangeArgs{Key: feedKey(userID), Start: 0, Stop: int64(limit - 1), Rev: true}).Result()
	if err != nil {
		return nil, err
	}
	items := make([]FeedItem, 0, len(members))
	for _, m := range members {
		var item FeedItem
		if err := json.Unmarshal([]byte(m), &item); err != nil {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func feedKey(userID uint64) string { return fmt.Sprintf("feed:%d", userID) }
