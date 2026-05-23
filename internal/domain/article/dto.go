package article

type CreateArticleRequest struct {
	Title    string                   `json:"title" binding:"required,max=255"`
	Slug     string                   `json:"slug" binding:"required,max=255"`
	Summary  string                   `json:"summary" binding:"required"`
	Category string                   `json:"category" binding:"required,max=100"`
	ReadTime string                   `json:"read_time" binding:"required,max=50"`
	Content  []map[string]interface{} `json:"content" binding:"required"`
}

type UpdateArticleRequest struct {
	Title    string                   `json:"title" binding:"required,max=255"`
	Slug     string                   `json:"slug" binding:"required,max=255"`
	Summary  string                   `json:"summary" binding:"required"`
	Category string                   `json:"category" binding:"required,max=100"`
	ReadTime string                   `json:"read_time" binding:"required,max=50"`
	Content  []map[string]interface{} `json:"content" binding:"required"`
}

type ArticleResponse struct {
	ID       uint                     `json:"id"`
	Slug     string                   `json:"slug"`
	Title    string                   `json:"title"`
	Summary  string                   `json:"summary"`
	Category string                   `json:"category"`
	Date     string                   `json:"date"`
	ReadTime string                   `json:"readTime"`
	Author   string                   `json:"author"`
	Content  []map[string]interface{} `json:"content"`
}

// Helper untuk mengubah Model DB "Article" menjadi Response DTO "ArticleResponse"
func ToArticleResponse(a Article) ArticleResponse {
	return ArticleResponse{
		ID:       a.ID,
		Slug:     a.Slug,
		Title:    a.Title,
		Summary:  a.Summary,
		Category: a.Category,
		Date:     a.CreatedAt.Format("Jan 02, 2006"), // Mengubah time.Time menjadi format "May 10, 2026"
		ReadTime: a.ReadTime,
		Author:   a.Author.Name,
		Content:  a.Content,
	}
}
