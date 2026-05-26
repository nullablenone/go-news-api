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

// Helper privat untuk mengambil data dari Redis
func (s *articleService) getCache(ctx context.Context, key string, dest interface{}) bool {
	cachedData, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(cachedData), dest); err != nil {
		return false
	}
	return true
}

// Helper privat untuk menyimpan ke Redis
func (s *articleService) setCache(ctx context.Context, key string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err == nil {
		s.rdb.Set(ctx, key, jsonData, cacheTTL)
	}
}

// Helper privat untuk membersihkan cache
func (s *articleService) clearArticleCache(slug string) {
	ctx := context.Background()
	s.rdb.Del(ctx, cacheKeyAllArticles)
	if slug != "" {
		s.rdb.Del(ctx, cacheKeyArticlePrefix+slug)
	}
}
