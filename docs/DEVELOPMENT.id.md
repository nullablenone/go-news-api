# Go News API - Dokumentasi Pengembang

Baca ini dalam [Bahasa Inggris](DEVELOPMENT.md)

Selamat datang di dokumentasi pengembang untuk **Go News API**. Dokumen ini ditulis khusus untuk memudahkan pengembangan di masa mendatang, proses onboarding pengembang baru, serta pemeliharaan kode. Dokumentasi ini berfokus pada detail teknis, arsitektur sistem, skema database, mekanisme caching, dan panduan untuk memperluas fitur aplikasi.

---

## Struktur Proyek & Navigasi

Basis kode ini menggunakan struktur Clean Architecture yang rapi. Berikut adalah direktori utama dan file entri penting:

- **`cmd/api/main.go`**: Bootstrapper aplikasi. Mengatur pemuatan konfigurasi, mode Gin, auto-migrasi database, inisialisasi client cache opsional, wiring dependency injection, dan menjalankan HTTP server.
- **`internal/shared/config/`**: Pemuatan konfigurasi:
  - `internal/shared/config/config.go`: Memuat file `.env` dan mengelompokkan konfigurasi ke `AppConfig`, `DBConfig`, `RedisConfig`, dan `JWTConfig`.
- **`internal/shared/database/`**: Inisialisasi driver database dan cache:
  - `internal/shared/database/postgresql.go`: Menyiapkan koneksi database PostgreSQL menggunakan GORM beserta konfigurasi pool-nya.
  - `internal/shared/database/redis.go`: Menginisialisasi koneksi client Redis dan melakukan verifikasi ping, serta mengembalikan `nil` jika tidak dikonfigurasi.
- **`internal/domain/`**: Logika bisnis inti yang dibagi per domain, berisi model skema, DTO, repositori, servis, dan handler:
  - **`internal/domain/user/`**: Registrasi, login, rotasi refresh token, model data pengguna, dan database seeder admin.
  - **`internal/domain/article/`**: Operasi CRUD artikel, pemetaan DTO, seeder artikel dummy, dan integrasi cache Redis opsional di `service.go`.
- **`internal/shared/middleware/`**: Interseptor HTTP Gin:
  - `internal/shared/middleware/auth.go`: Mengimplementasikan verifikasi JWT, middleware refresh token, dan penegakan peran (role verification).
- **`internal/shared/utils/`**: Fungsi pembantu umum:
  - `internal/shared/utils/jwt.go`: Berfungsi untuk pembuatan dan validasi access/refresh token JWT.
- **`internal/routes/`**: Konfigurasi routing HTTP:
  - `internal/routes/routes.go`: Menghubungkan endpoint API ke handler terkait dan mendefinisikan aturan CORS.
- **`docs/`**: Dokumentasi eksternal dan koleksi pengujian API:
  - `docs/go-news-api.postman_collection.json`: Berkas koleksi Postman untuk pengujian API.

---

## Arsitektur Teknis & Aliran Permintaan (Request Flow)

Aplikasi ini menggunakan pola **Clean Architecture** untuk memisahkan tanggung jawab (separation of concerns) ke dalam 4 lapisan:

1. **Model Domain & DTO**: Struktur data murni yang mendefinisikan skema database dan format permintaan/respons API.
2. **Lapisan Repositori (Repository)**: Membungkus operasi database menggunakan GORM, mengisolasi kueri SQL/engine DB dari modul lainnya.
3. **Lapisan Layanan (Service)**: Mengimplementasikan logika bisnis inti aplikasi serta kebijakan manipulasi data dan caching.
4. **Lapisan Handler (Controller)**: Menangani ekstraksi permintaan HTTP, validasi payload JSON menggunakan binding tag Gin, dan pengembalian kode status HTTP beserta format respons.

### Siklus Hidup Permintaan HTTP (End-to-End)

Agar sistem lebih mudah dipahami oleh pengembang, siklus hidup permintaan HTTP dibagi menjadi dua alur yang disederhanakan:

#### 1. Aliran Permintaan Baca Publik (Pola Cache-Aside)
Alur ini digunakan saat mengambil artikel. Jika data sudah di-cache di Redis, respons akan segera dikembalikan untuk menghindari kueri langsung ke database PostgreSQL.

```mermaid
sequenceDiagram
    actor Client as Client (Reader)
    participant API as API Handler
    participant Cache as Redis
    participant DB as PostgreSQL

    Client->>API: GET /public/articles
    API->>Cache: Periksa Cache (articles:all)
    alt Cache Hit
        Cache-->>Client: Kembalikan Artikel dari Cache
    else Cache Miss
        API->>DB: Kueri Artikel (Memuat Relasi Author)
        DB-->>API: Kembalikan Artikel
        API->>Cache: Simpan ke Cache (TTL 1 Jam)
        API-->>Client: Kembalikan Artikel
    end
```

#### 2. Aliran Permintaan Tulis Admin Terautentikasi (Otorisasi & Invalidasi Cache)
Alur ini dipicu saat melakukan mutasi data artikel. Alur ini memvalidasi kredensial pengguna dan secara otomatis menghapus data cache yang usang.

```mermaid
sequenceDiagram
    actor Admin as Admin (CMS)
    participant Mid as Middleware (JWT & Role)
    participant API as API Handler
    participant DB as PostgreSQL
    participant Cache as Redis

         Admin->>Mid: POST /admin/article (dengan JWT Bearer Token)
         alt Token Tidak Valid atau Peran != admin
             Mid-->>Admin: Kembalikan 401 Unauthorized / 403 Forbidden
         else Token Valid & Peran == admin
             Mid->>API: Teruskan Permintaan Terotorisasi
             API->>DB: Simpan/Perbarui/Hapus Artikel
             DB-->>API: DB Sukses (ID Dibuat/Diperbarui)
             API->>Cache: Bersihkan Cache (Kunci articles:all & slug terkait)
             API-->>Admin: Kembalikan Respons Sukses (201 Created / 200 OK)
         end
```

---

## Skema Database & Serialisasi Kustom

Tabel relasional dikelola secara otomatis menggunakan fitur **GORM AutoMigrate** di dalam berkas `cmd/api/main.go`.

### 1. Entitas User (`internal/domain/user/model.go`)
Menyimpan informasi pengguna/admin. Mendukung fitur penghapusan logis (soft delete):
- `ID` (Primary Key, uint)
- `Name` (varchar(100), tidak null)
- `Email` (varchar(150), unik, indeks, tidak null)
- `Password` (varchar(255), tidak null) - Disimpan dalam bentuk hash Bcrypt.
- `Role` (varchar(50), default 'user', tidak null) - Menentukan tingkat otorisasi RBAC (`admin` atau `user`).
- `CreatedAt`, `UpdatedAt` (timestamptz)
- `DeletedAt` (indeks timestamptz) - GORM soft delete support.

### 2. Entitas Article (`internal/domain/article/model.go`)
Menyimpan data dokumen artikel berita.
- `ID` (Primary Key, uint)
- `Title` (varchar(255), tidak null)
- `Slug` (varchar(255), unik, indeks, tidak null)
- `Summary` (text, tidak null)
- `Category` (varchar(100), tidak null)
- `ReadTime` (varchar(50), tidak null)
- `AuthorID` (uint, foreign key mengarah ke `users(id)`)
- `Content` (`[]map[string]interface{}`, tipe kolom `jsonb`, terindeks GIN)
- `CreatedAt`, `UpdatedAt` (timestamptz)
- `DeletedAt` (indeks timestamptz)

#### Serialisasi Konten Dinamis (Dynamic Content Serialization)
Salah satu fitur teknis utama dari modul artikel ini terletak pada kolom `Content`:
```go
Content []map[string]interface{} `gorm:"serializer:json;type:jsonb;index:,type:gin;not null"`
```
- **Proses Serialisasi**: GORM secara otomatis melakukan *marshal* (konversi dari Go slice ke teks string JSON) saat operasi tulis, dan *unmarshal* (kembali menjadi Go data structure `[]map[string]interface{}`) saat operasi baca.
- **Tipe Data Database**: Kolom menggunakan tipe data asli PostgreSQL `jsonb` yang sangat optimal dalam penyimpanan data JSON biner. Fitur ini memungkinkan penyimpanan data blok editor modern (seperti Editor.js atau Quill editor).
- **Pengindeksan**: Dibuat indeks **GIN** (Generalized Inverted Index) di PostgreSQL pada kolom `Content`, memastikan kueri pencarian data bersarang (nested JSON) tetap sangat cepat.

---

## Mekanisme Caching Redis & Invalidasi Cache

Untuk menahan beban kueri saat lalu lintas padat, endpoint publik mengimplementasikan pola **Cache-Aside (Lazy Loading)** secara opsional. Jika Redis tidak dikonfigurasi atau tidak tersedia saat startup, aplikasi akan langsung menggunakan PostgreSQL sebagai fallback.

### Konfigurasi Cache
- **Masa Aktif (TTL)**: 1 Jam (`cacheTTL` di dalam berkas `internal/domain/article/service.go`).
- **Kebijakan Kunci (Keys)**:
  - `articles:all`: Menyimpan seluruh daftar artikel yang ditayangkan ke publik.
  - `article:slug:<slug>`: Menyimpan detail artikel tunggal berdasarkan slug uniknya.

### Kebijakan Invalidasi Cache
Untuk mencegah pembacaan data basi (stale data), operasi penulisan data (mutasi) wajib menghapus cache yang terpengaruh. Berikut adalah aturan penghapusan cache yang diimplementasikan di `internal/domain/article/service.go`:

| Operasi Tulis | Endpoint Terkait | Kunci Cache yang Dihapus | Alasan & Rationale |
|---|---|---|---|
| **Create (Buat)** | `POST /admin/article` | `articles:all` | Daftar artikel di publik telah berubah dan harus dimuat ulang dari database pada request berikutnya. |
| **Update (Perbarui)** | `PUT /admin/article/:id` | `articles:all`<br>`article:slug:<old_slug>`<br>`article:slug:<new_slug>` | Daftar artikel berubah. Jika pengembang/admin mengubah slug artikel, sistem akan membersihkan cache dari slug lama dan slug baru agar tidak terjadi eror 404 cache. |
| **Delete (Hapus)** | `DELETE /admin/article/:id` | `articles:all`<br>`article:slug:<slug>` | Daftar publik berubah; detail cache dari artikel yang dihapus segera disingkirkan agar tidak bisa diakses kembali. |

---

## Autentikasi, Otorisasi, & RBAC

API ini menggunakan mekanisme **Stateless JWT Tokens** untuk keamanan rute.

### 1. Struktur Token
- Token dibuat melalui berkas pembantu `internal/shared/utils/jwt.go`.
- Menyimpan klaim kustom: `UserID`, `Role`, dan `TokenType`.
- Ditandatangani menggunakan algoritma enkripsi simetris `HS256` dengan kunci rahasia yang diambil dari variabel environment `JWT_SECRET`.
- Masa kedaluwarsa token: 40 menit untuk access token, 7 hari untuk refresh token.

### 2. Aliran Interseptor Rute
- **`AuthMiddleware` (`internal/shared/middleware/auth.go`)**: Mengadang permintaan ke rute terproteksi, membaca header `Authorization` dengan format `Bearer <token>`, memverifikasi keabsahan access JWT, dan mengekstrak klaim user. Jika token valid, informasi `user_id` dan `role` disimpan ke dalam konteks permintaan Gin (`c.Set("user_id", ...)` dan `c.Set("role", ...)`).
- **`RoleMiddleware` (`internal/shared/middleware/auth.go`)**: Memeriksa data `role` yang tersimpan pada konteks Gin. Membandingkannya dengan daftar peran yang diperbolehkan mengakses rute (contoh: `RoleMiddleware("admin")`).

---

## Konfigurasi Lingkungan (Environment) & Menjalankan Aplikasi

### 1. Konfigurasi Variabel (`.env`)
Salin file `.env.example` menjadi `.env` di direktori utama proyek, kemudian sesuaikan nilainya:
```ini
APP_NAME=go-news-api
APP_ENV=development
APP_PORT=8888

DB_HOST=localhost
DB_USER=postgres
DB_PASS=isi_password_database_anda
DB_NAME=news_db
DB_PORT=5432
DB_SSLMODE=disable
DB_MAX_OPEN_CONNS=50
DB_MAX_IDLE_CONNS=10

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

JWT_SECRET=buat_kunci_rahasia_jwt_yang_panjang_dan_acak
```

### 2. Menjalankan Dependensi Eksternal (Docker)
Jika Anda belum memiliki instance PostgreSQL atau Redis yang berjalan di sistem lokal, gunakan perintah Docker berikut (Redis bersifat opsional):
```bash
# Menjalankan Redis (Opsional)
docker run --name news-api-redis -p 6379:6379 -d redis

# Menjalankan PostgreSQL
docker run --name news-api-postgres -e POSTGRES_DB=news_db -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=isi_password_database_anda -p 5432:5432 -d postgres
```

### 3. Perintah Eksekusi
Gunakan perintah Go berikut untuk memproses dependensi dan menjalankan server:
```bash
# Mengunduh dependensi
go mod tidy

# Menjalankan server aplikasi
go run ./cmd/api
```
Secara default, HTTP server akan berjalan pada port `8888` (dapat diubah melalui `APP_PORT`). Sinkronisasi tabel database (auto-migration) berjalan otomatis saat inisialisasi aplikasi.

### 4. Akun Admin Bawaan (Admin Seeder)
Saat pertama kali dijalankan, `cmd/api/main.go` akan memanggil `RunAdminSeeder` untuk memeriksa keberadaan user dengan role `"admin"`. Jika kosong, sistem otomatis membuat satu akun admin default:
- **Email**: `admin@news.com`
- **Password**: `#admin123` *(Penting: Pastikan untuk menyertakan awalan karakter pagar `#` sesuai implementasi pada kode seeder)*

---

## Pengujian API dengan Postman

Tersedia berkas koleksi Postman untuk mempercepat pengujian setiap rute API.

- **Lokasi Berkas**: `docs/go-news-api.postman_collection.json`
- **Cara Mengimpor**:
  1. Buka aplikasi client Postman.
  2. Klik tombol **Import** di bagian panel navigasi kiri atas.
  3. Pilih dan unggah berkas `go-news-api.postman_collection.json` dari repositori ini.
  4. Seluruh endpoint API (Registrasi, Login, Rute Publik, dan Rute Manajemen Admin) beserta contoh payload-nya langsung siap digunakan.

---

## Panduan Pengembangan: Cara Memperluas Aplikasi

Gunakan langkah-langkah terstruktur di bawah ini sebagai acuan (cookbook) untuk menambahkan modul atau domain baru (contoh: Menambahkan fitur **Komentar (Comments)** ke artikel):

### Langkah 1: Buat Direktori Domain Baru
Buat folder baru di bawah `internal/domain/`:
```bash
mkdir internal/domain/comment
```

### Langkah 2: Definisikan Model Database
Buat file `internal/domain/comment/model.go` berisi definisi skema GORM:
```go
package comment

import (
	"time"
	"gorm.io/gorm"
)

type Comment struct {
	ID        uint   `gorm:"primaryKey"`
	ArticleID uint   `gorm:"not null"`
	Author    string `gorm:"size:100;not null"`
	Body      string `gorm:"type:text;not null"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

### Langkah 3: Definisikan Request & Response DTO
Buat file `internal/domain/comment/dto.go` untuk menangani struktur input validasi dan format keluaran API:
```go
package comment

type CreateCommentRequest struct {
	Author string `json:"author" binding:"required,max=100"`
	Body   string `json:"body" binding:"required"`
}

type CommentResponse struct {
	ID        uint   `json:"id"`
	Author    string `json:"author"`
	Body      string `json:"body"`
	CreatedAt string `json:"createdAt"`
}
```

### Langkah 4: Tulis Kontrak Repositori & Implementasinya
Buat file `internal/domain/comment/repo.go` untuk mengisolasi query database GORM:
```go
package comment

import "gorm.io/gorm"

type CommentRepository interface {
	Create(comment *Comment) error
	FindByArticleID(articleID uint) ([]Comment, error)
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(c *Comment) error {
	return r.db.Create(c).Error
}

func (r *commentRepository) FindByArticleID(articleID uint) ([]Comment, error) {
	var comments []Comment
	err := r.db.Where("article_id = ?", articleID).Find(&comments).Error
	return comments, err
}
```

### Langkah 5: Implementasikan Lapisan Bisnis (Service)
Buat file `internal/domain/comment/service.go`. Di bagian ini Anda bisa menambahkan logika validasi bisnis tambahan atau integrasi cache jika diperlukan:
```go
package comment

type CommentService interface {
	AddComment(articleID uint, req CreateCommentRequest) (Comment, error)
	GetCommentsByArticle(articleID uint) ([]Comment, error)
}

type commentService struct {
	repo CommentRepository
}

func NewCommentService(repo CommentRepository) CommentService {
	return &commentService{repo: repo}
}

// Implementasikan method di bawah ini...
```

### Langkah 6: Implementasikan HTTP Handler (Controller)
Buat file `internal/domain/comment/handler.go` untuk memproses binding payload JSON dan respons HTTP:
```go
package comment

import (
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	service CommentService
}

func NewCommentHandler(service CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

func (h *CommentHandler) AddComment(c *gin.Context) {
    // 1. Ambil parameter dari rute/kueri
    // 2. Lakukan c.ShouldBindJSON()
    // 3. Panggil method service
    // 4. Kirim respons c.JSON()
}
```

### Langkah 7: Daftarkan Skema & Hubungkan Dependensi (Wiring)
Buka berkas **`cmd/api/main.go`** dan tambahkan penyesuaian berikut:
1. Daftarkan model baru pada fungsi `AutoMigrate` GORM:
   ```go
   db.AutoMigrate(article.Article{}, user.User{}, comment.Comment{})
   ```
2. Hubungkan dependency injection instans repositori, service, dan handler:
   ```go
   commentRepo := comment.NewCommentRepository(db)
   commentService := comment.NewCommentService(commentRepo)
   commentHandler := comment.NewCommentHandler(commentService)
   ```
3. Oper instans handler baru ke dalam fungsi inisialisasi rute:
   ```go
   router := routes.SetRoutes(cfg, articleHandler, userHandler, commentHandler)
   ```

### Langkah 8: Konfigurasikan Endpoint Rute API
Buka berkas **`internal/routes/routes.go`** untuk mendaftarkan endpoint HTTP rute:
- Contoh mendaftarkan sebagai rute publik pembaca:
  ```go
  publicArticleRoutes.POST("/articles/:id/comments", commentHandler.AddComment)
  ```
