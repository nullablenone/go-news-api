# Go News API

Baca ini dalam [Bahasa Inggris](README.md) 🇺🇸

Go News API adalah backend RESTful API berkinerja tinggi yang dirancang untuk portal berita, blog. Dibangun dengan **Golang**, proyek ini menerapkan prinsip **Clean Architecture** untuk memastikan skalabilitas dan kemudahan pemeliharaan kode.

## Fitur Utama
* **Arsitektur Bersih (Clean Architecture)**: Pemisahan logika bisnis yang jelas antara *Domain*, *Repository*, *Service*, dan *Handler*.
* **Penyimpanan Konten Dinamis**: Menggunakan tipe data `[]map[string]interface{}` dengan kolom `jsonb` PostgreSQL untuk menyimpan struktur data dinamis dari berbagai *rich text editor* di frontend (seperti Editor.js, Quill, dsb).
* **Caching Super Cepat dengan Redis**: Menerapkan pola *Cache-Aside* untuk *endpoint* publik, sehingga meminimalkan beban database saat terjadi lonjakan pengunjung. Dilengkapi dengan fitur Invalidasi Cache otomatis pada operasi CRUD admin.
* **Autentikasi & Otorisasi**: Dilindungi oleh Token JWT dan *Role-Based Access Control* (RBAC).
* **Pemisahan Rute**: Rute Publik untuk pembaca (Website A) dan Rute Admin Private untuk dasbor CMS (Website B).
* **Siap untuk Frontend**: Telah dikonfigurasi dengan CORS yang aman untuk diintegrasikan dengan framework seperti React, Vue, atau Angular.

## Tech Stack
* **Bahasa**: Go (1.25.0)
* **Framework Web**: [Gin-Gonic](https://gin-gonic.com/)
* **Database**: PostgreSQL (via [GORM](https://gorm.io/))
* **Cache**: Redis (via `go-redis/v9`)
* **Keamanan**: JWT (`golang-jwt/v5`) & Bcrypt

## Struktur Proyek
```text
.
├── config/              # Konfigurasi Env, Koneksi Database (PostgreSQL & Redis)
├── internal/
│   ├── domain/          # Entitas inti, DTO, Repository, Service, dan Handler
│   │   ├── article/     # Modul Artikel (termasuk helper Cache)
│   │   └── user/        # Modul User (Autentikasi & RBAC)
│   ├── middleware/      # Middleware JWT & Validasi Role
│   └── utils/           # Fungsi bantuan (Generator JWT, dll)
├── routes/              # Konfigurasi routing Gin & CORS
├── main.go              # Entry point aplikasi
└── .env                 # Variabel environment
````

## Cara Menjalankan 

### 1. Prasyarat

- [Go](https://golang.org/dl/) terinstal
    
- [PostgreSQL](https://www.postgresql.org/) berjalan
    
- [Docker](https://www.docker.com/) (untuk menjalankan Redis)
    

### 2. Konfigurasi Environment

Salin file `.env.example` menjadi `.env` dan sesuaikan kredensialnya:

Cuplikan kode

```
DB_HOST=localhost
DB_USER=postgres
DB_PASS=password_kamu
DB_NAME=news_db
DB_PORT=5432
DB_SSLMODE=disable

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

JWT_SECRET=secret
```

### 3. Jalankan Redis via Docker

Bash

```
docker run --name news-api-redis -p 6379:6379 -d redis
```

### 4. Instal Dependensi & Jalankan

Bash

```
go mod tidy
go run main.go
```

Aplikasi akan berjalan di `http://localhost:8888`. Database akan di-migrasi secara otomatis (_AutoMigrate_).

## Dokumentasi API (Ringkasan)

**Autentikasi (Public):**
* `POST /register` - Pendaftaran akun pengguna baru
* `POST /login` - Autentikasi pengguna dan mendapatkan token JWT

**Publik / Website Pembaca (Tanpa Token - Pola Cache-Aside):**
* `GET /public/articles` - Mengambil daftar semua artikel untuk pembaca
* `GET /public/articles/:slug` - Mengambil detail artikel berdasarkan slug

**Admin / Dashboard CMS (Membutuhkan Bearer Token, Role: 'admin'):**
* `GET /admin/articles` - Mengambil semua artikel (untuk tabel manajemen data)
* `GET /admin/article/:id` - Mengambil detail spesifik artikel berdasarkan ID
* `POST /admin/article` - Membuat artikel baru (Otomatis menghapus cache publik)
* `PUT /admin/article/:id` - Memperbarui artikel yang sudah ada (Otomatis menghapus cache publik)
* `DELETE /admin/article/:id` - Menghapus artikel (Otomatis menghapus cache publik)