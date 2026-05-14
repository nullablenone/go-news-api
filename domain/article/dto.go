package article

type CreateArticleRequest struct {
	Title   string `json:"title" binding:"required,max=225"`
	Slug    string `json:"slug" binding:"required,max=225"`
	Content string `json:"content" binding:"required"`
}
