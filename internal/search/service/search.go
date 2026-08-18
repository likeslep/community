// Package service 是 search-service 的业务逻辑层。
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

// IndexName 是统一搜索索引名。
const IndexName = "community"

// Document 是搜索文档。
type Document struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"` // article / question / answer
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	AuthorID uint64   `json:"author_id"`
	Tags     []string `json:"tags"`
}

// SearchResult 是搜索结果。
type SearchResult struct {
	ID       string
	Type     string
	Title    string
	Snippet  string
	AuthorID uint64
}

// SearchService 是搜索业务逻辑层，Elasticsearch 为 Search Index（MySQL 为 Source of Truth）。
type SearchService struct {
	es *elasticsearch.Client
}

// NewSearchService 构造。
func NewSearchService(es *elasticsearch.Client) *SearchService {
	return &SearchService{es: es}
}

// EnsureIndex 创建索引（含 mapping）。
func (s *SearchService) EnsureIndex(ctx context.Context) error {
	mapping := `{
		"mappings": {
			"properties": {
				"id": {"type": "keyword"},
				"type": {"type": "keyword"},
				"title": {"type": "text"},
				"content": {"type": "text"},
				"author_id": {"type": "long"},
				"tags": {"type": "keyword"}
			}
		}
	}`
	res, err := s.es.Indices.Create(IndexName, s.es.Indices.Create.WithContext(ctx), s.es.Indices.Create.WithBody(strings.NewReader(mapping)))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

// Index 写入文档。
func (s *SearchService) Index(ctx context.Context, doc Document) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	res, err := s.es.Index(IndexName, bytes.NewReader(body), s.es.Index.WithContext(ctx), s.es.Index.WithDocumentID(doc.ID))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

// Delete 删除文档。
func (s *SearchService) Delete(ctx context.Context, id string) error {
	res, err := s.es.Delete(IndexName, id, s.es.Delete.WithContext(ctx))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

// Search 执行搜索（关键词 + 过滤 + 高亮 + 分页）。
func (s *SearchService) Search(ctx context.Context, keyword, typeFilter string, tags []string, page, pageSize int) ([]SearchResult, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	boolQuery := map[string]any{
		"must": []any{
			map[string]any{
				"multi_match": map[string]any{"query": keyword, "fields": []string{"title", "content"}},
			},
		},
	}
	var filters []any
	if typeFilter != "" {
		filters = append(filters, map[string]any{"term": map[string]any{"type": typeFilter}})
	}
	if len(tags) > 0 {
		filters = append(filters, map[string]any{"terms": map[string]any{"tags": tags}})
	}
	if len(filters) > 0 {
		boolQuery["filter"] = filters
	}

	query := map[string]any{
		"query": map[string]any{"bool": boolQuery},
		"highlight": map[string]any{
			"fields": map[string]any{"title": map[string]any{}, "content": map[string]any{}},
		},
		"from": (page - 1) * pageSize,
		"size": pageSize,
	}
	body, err := json.Marshal(query)
	if err != nil {
		return nil, 0, err
	}

	res, err := s.es.Search(
		s.es.Search.WithContext(ctx),
		s.es.Search.WithIndex(IndexName),
		s.es.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	var decoded struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source    Document            `json:"_source"`
				Highlight map[string][]string `json:"highlight"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&decoded); err != nil {
		return nil, 0, err
	}

	results := make([]SearchResult, 0, len(decoded.Hits.Hits))
	for _, h := range decoded.Hits.Hits {
		snippet := ""
		if v, ok := h.Highlight["content"]; ok && len(v) > 0 {
			snippet = v[0]
		}
		results = append(results, SearchResult{
			ID: h.Source.ID, Type: h.Source.Type, Title: h.Source.Title,
			Snippet: snippet, AuthorID: h.Source.AuthorID,
		})
	}
	return results, decoded.Hits.Total.Value, nil
}
