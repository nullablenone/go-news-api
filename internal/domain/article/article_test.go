package article

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockArticleRepo struct {
	items  map[uint]Article
	bySlug map[string]uint
	nextID uint
}

func newMockArticleRepo() *mockArticleRepo {
	return &mockArticleRepo{
		items:  make(map[uint]Article),
		bySlug: make(map[string]uint),
		nextID: 1,
	}
}

func (m *mockArticleRepo) Create(a *Article) error {
	a.ID = m.nextID
	m.nextID++
	m.items[a.ID] = *a
	m.bySlug[a.Slug] = a.ID
	return nil
}

func (m *mockArticleRepo) FindAll() ([]Article, error) {
	out := make([]Article, 0, len(m.items))
	for _, a := range m.items {
		out = append(out, a)
	}
	return out, nil
}

func (m *mockArticleRepo) FindByID(id uint) (Article, error) {
	a, ok := m.items[id]
	if !ok {
		return Article{}, errors.New("not found")
	}
	return a, nil
}

func (m *mockArticleRepo) FindBySlug(slug string) (Article, error) {
	id, ok := m.bySlug[slug]
	if !ok {
		return Article{}, errors.New("not found")
	}
	return m.items[id], nil
}

func (m *mockArticleRepo) Update(a *Article) error {
	old, ok := m.items[a.ID]
	if !ok {
		return errors.New("not found")
	}
	delete(m.bySlug, old.Slug)
	m.items[a.ID] = *a
	m.bySlug[a.Slug] = a.ID
	return nil
}

func (m *mockArticleRepo) Delete(id uint) error {
	a, ok := m.items[id]
	if !ok {
		return errors.New("not found")
	}
	delete(m.bySlug, a.Slug)
	delete(m.items, id)
	return nil
}

func TestArticleServiceWithoutRedis(t *testing.T) {
	repo := newMockArticleRepo()
	svc := NewArticleService(repo, nil)

	created, err := svc.CreateArticle(1, CreateArticleRequest{
		Title:    "Go Testing",
		Slug:     "go-testing",
		Summary:  "Summary",
		Category: "Tech",
		ReadTime: "3 min",
		Content:  []map[string]interface{}{{"type": "paragraph", "text": "Hello"}},
	})
	require.NoError(t, err)
	assert.Equal(t, uint(1), created.ID)

	all, err := svc.GetAllArticles()
	require.NoError(t, err)
	assert.Len(t, all, 1)

	bySlug, err := svc.GetArticleBySlug("go-testing")
	require.NoError(t, err)
	assert.Equal(t, "Go Testing", bySlug.Title)

	updated, err := svc.UpdateArticle(created.ID, UpdateArticleRequest{
		Title:    "Go Testing Updated",
		Slug:     "go-testing-updated",
		Summary:  "Updated Summary",
		Category: "Tech",
		ReadTime: "4 min",
		Content:  []map[string]interface{}{{"type": "paragraph", "text": "Updated"}},
	})
	require.NoError(t, err)
	assert.Equal(t, "go-testing-updated", updated.Slug)

	err = svc.DeleteArticle(created.ID)
	require.NoError(t, err)

	_, err = svc.GetArticleByID(created.ID)
	assert.Error(t, err)
}
