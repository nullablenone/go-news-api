package article

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func RunArticleSeeder(db *gorm.DB) {
	var count int64
	db.Model(&Article{}).Count(&count)

	if count < 2000 {
		log.Println("====== [SEEDER] Starting insert of 2000 dummy articles... ======")

		var articles []Article
		for i := 1; i <= 2000; i++ {
			articles = append(articles, Article{
				Title:    fmt.Sprintf("Dummy Article #%d for Load Test", i),
				Slug:     fmt.Sprintf("dummy-article-%d", i),
				Summary:  "This is a dummy article summary to test Redis cache performance.",
				Category: "Benchmark",
				ReadTime: "3 min read",
				AuthorID: 1,
				Content: []map[string]interface{}{
					{"type": "paragraph", "text": "This is dummy article content to load-test database queries."},
				},
			})
		}

		if err := db.Create(&articles).Error; err != nil {
			log.Printf("[SEEDER ERROR] failed to insert dummy articles: %v\n", err)
		} else {
			log.Println("====== [SEEDER] Successfully inserted 2000 dummy articles! ======")
		}
	} else {
		log.Println("====== [SEEDER] 2000 dummy articles already exist ======")
	}
}
