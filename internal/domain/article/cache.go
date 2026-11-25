package article

import (
	"context"
	"encoding/json"
	"time"
)

const (
	cacheTTL              = 1 * time.Hour
	cacheKeyAllArticles   = "articles:all"
	cacheKeyArticlePrefix = "article:slug:"
)

func (s *articleService) getCache(ctx context.Context, key string, dest interface{}) bool {
	if s.rdb == nil {
		return false
	}
	cachedData, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(cachedData), dest); err != nil {
		return false
	}
	return true
}

func (s *articleService) setCache(ctx context.Context, key string, data interface{}) {
	if s.rdb == nil {
		return
	}
	jsonData, err := json.Marshal(data)
	if err == nil {
		s.rdb.Set(ctx, key, jsonData, cacheTTL)
	}
}

func (s *articleService) clearArticleCache(slug string) {
	if s.rdb == nil {
		return
	}
	ctx := context.Background()
	s.rdb.Del(ctx, cacheKeyAllArticles)
	if slug != "" {
		s.rdb.Del(ctx, cacheKeyArticlePrefix+slug)
	}
}
