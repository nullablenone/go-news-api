# Go News API

Baca ini dalam [Bahasa Inggris](README.md) 🇺🇸

Go News API adalah backend RESTful API berkinerja tinggi yang dirancang untuk portal berita, blog. Dibangun dengan **Golang**, proyek ini menerapkan prinsip **Clean Architecture** untuk memastikan skalabilitas dan kemudahan pemeliharaan kode.

## Fitur Utama

- **Arsitektur Bersih (Clean Architecture)**: Pemisahan logika bisnis yang jelas antara _Domain_, _Repository_, _Service_, dan _Handler_.
- **Penyimpanan Konten Dinamis**: Menggunakan tipe data `[]map[string]interface{}` dengan kolom `jsonb` PostgreSQL untuk menyimpan struktur data dinamis dari berbagai _rich text editor_ di frontend (seperti Editor.js, Quill, dsb).
- **Caching Super Cepat dengan Redis**: Menerapkan pola _Cache-Aside_ untuk _endpoint_ publik, sehingga meminimalkan beban database saat terjadi lonjakan pengunjung. Dilengkapi dengan fitur Invalidasi Cache otomatis pada operasi CRUD admin.
- **Autentikasi & Otorisasi**: Dilindungi oleh Token JWT dan _Role-Based Access Control_ (RBAC).
- **Pemisahan Rute**: Rute Publik untuk pembaca (Website A) dan Rute Admin Private untuk dasbor CMS (Website B).
- **Siap untuk Frontend**: Telah dikonfigurasi dengan CORS yang aman untuk diintegrasikan dengan framework seperti React, Vue, atau Angular.

## Tech Stack

- **Bahasa**: Go (1.25.0)
- **Framework Web**: [Gin-Gonic](https://gin-gonic.com/)
- **Database**: PostgreSQL (via [GORM](https://gorm.io/))
- **Cache**: Redis (via `go-redis/v9`)
- **Keamanan**: JWT (`golang-jwt/v5`) & Bcrypt

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
```

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

- `POST /register` - Pendaftaran akun pengguna baru
- `POST /login` - Autentikasi pengguna dan mendapatkan token JWT

**Publik / Website Pembaca (Tanpa Token - Pola Cache-Aside):**

- `GET /public/articles` - Mengambil daftar semua artikel untuk pembaca
- `GET /public/articles/:slug` - Mengambil detail artikel berdasarkan slug

**Admin / Dashboard CMS (Membutuhkan Bearer Token, Role: 'admin'):**

- `GET /admin/articles` - Mengambil semua artikel (untuk tabel manajemen data)
- `GET /admin/article/:id` - Mengambil detail spesifik artikel berdasarkan ID
- `POST /admin/article` - Membuat artikel baru (Otomatis menghapus cache publik)
- `PUT /admin/article/:id` - Memperbarui artikel yang sudah ada (Otomatis menghapus cache publik)
- `DELETE /admin/article/:id` - Menghapus artikel (Otomatis menghapus cache publik)

## Kredensial Admin Bawaan (Default Admin)

Aplikasi ini dilengkapi dengan fitur _auto-seeder_ modular. Jika database berjalan dalam keadaan kosong, sistem akan otomatis membuat satu akun admin. Gunakan akun berikut untuk login dan menguji endpoint `/admin/*`:

- **Email:** `admin@news.com`
- **Password:** `admin123`

## Pengujian API dengan Postman

Untuk mempermudah pengujian endpoint, repositori ini telah dilengkapi dengan file Postman Collection.

**Cara Penggunaan:**

1. Buka aplikasi [Postman](https://www.postman.com/).
2. Klik tombol **Import** di sudut kiri atas.
3. Unggah file `docs/go-news-api.postman_collection.json` dari repositori ini.
4. Seluruh _endpoint_ API (Public & Admin) beserta contoh _payload_ JSON-nya sudah tersedia dan siap dijalankan.
