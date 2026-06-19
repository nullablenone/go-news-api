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
		log.Println("====== [SEEDER] Memulai proses insert 2000 artikel dummy... ======")

		var articles []Article
		for i := 1; i <= 2000; i++ {
			articles = append(articles, Article{
				Title:    fmt.Sprintf("Artikel Dummy Ke-%d untuk Load Test", i),
				Slug:     fmt.Sprintf("artikel-dummy-ke-%d", i),
				Summary:  "Ini adalah ringkasan artikel dummy untuk menguji performa Redis Cache.",
				Category: "Benchmark",
				ReadTime: "3 min read",
				AuthorID: 1, // ID 1 adalah milik admin@news.com dari seeder user
				Content: []map[string]interface{}{
					{"type": "paragraph", "text": "Ini adalah konten artikel dummy untuk membebani query database."},
				},
			})
		}

		// GORM akan melakukan Batch Insert sehingga prosesnya instan
		if err := db.Create(&articles).Error; err != nil {
			log.Printf("[SEEDER ERROR] Gagal insert artikel dummy: %v\n", err)
		} else {
			log.Println("====== [SEEDER] Berhasil memasukkan 2000 artikel dummy! ======")
		}
	} else {
		log.Println("====== [SEEDER] 2000 Artikel dummy sudah tersedia ======")
	}
}
