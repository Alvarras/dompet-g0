# 🚀 Go Budget API

Go Budget adalah aplikasi manajemen anggaran (expense tracker) yang dibangun dengan Go (Golang) menggunakan framework Echo. Aplikasi ini menyediakan API untuk mengelola budget dan pengeluaran dengan fitur autentikasi JWT.

## Fitur

- Autentikasi JWT
- Manajemen Budget
- Manajemen Pengeluaran
- Tracking Penggunaan Budget

## 🛠️ Teknologi

- Go 1.24+
- Echo Framework
- MySQL
- JWT
- GORM
- Go Validator

# Project pattern and layering
Proyek ini dibangun menggunakan bahasa Go dengan struktur direktori yang terorganisir untuk memfasilitasi pengembangan dan pemeliharaan. Secara umum, kode aplikasi utama terletak di dalam direktori `internal` untuk memastikan bahwa paket-paket di dalamnya tidak dapat diimpor oleh proyek eksternal, mengikuti konvensi Go. Titik masuk aplikasi didefinisikan dalam direktori `cmd`.

Arsitektur aplikasi ini mengikuti pola perancangan layering pattern yang clean, memisahkan setiap bagian berdasarkan fungsinya secara efektif:

* **`cmd/`**: Berisi program utama aplikasi.
    * `main.go`: Merupakan titik masuk (entry point) aplikasi, bertanggung jawab untuk inisialisasi dan menjalankan server.
* **`internal/`**: Berisi seluruh kode inti aplikasi yang tidak dimaksudkan untuk diimpor oleh aplikasi atau pustaka eksternal. Di dalamnya, kode diorganisir lebih lanjut ke dalam lapisan-lapisan berikut:

    * **`config/`**: Manajemen konfigurasi aplikasi.
        * `config.go`: Logika untuk memuat dan mengelola konfigurasi.
    * **`controllers/`**: **(Lapisan Presentasi)** Menangani permintaan HTTP, melakukan validasi input, dan mengembalikan respons.
        * `auth_controller.go`: Endpoint untuk autentikasi pengguna.
        * `budget_controller.go`: Endpoint untuk manajemen anggaran.
        * `expense_controller.go`: Endpoint untuk pencatatan pengeluaran.
    * **`db/`**: **(Lapisan Infrastruktur)** Penanganan koneksi dan interaksi dengan database.
        * `db.go`: Inisialisasi dan manajemen koneksi database.
    * **`dtos/`**: **(Lapisan Logika Bisnis/Presentasi)** Data Transfer Objects (Objek Transfer Data) yang digunakan untuk komunikasi antar lapisan.
        * `requests/`: Struktur untuk validasi data masukan (input).
            * `auth_request.go`, `budget_request.go`, `expense_request.go`
        * `responses/`: Struktur untuk format data keluaran (output).
            * `auth_response.go`, `budget_response.go`, `common_response.go`, `expense_response.go`
    * **`middlewares/`**: **(Lapisan Presentasi)** Menangani aspek lintas-fungsi seperti validasi autentikasi sebelum permintaan mencapai controller.
        * `auth_middleware.go`: Middleware untuk autentikasi.
    * **`models/`**: **(Lapisan Akses Data/Domain)** Entitas domain yang merepresentasikan struktur data atau tabel-tabel dalam database.
        * `budget.go`, `expense.go`, `user.go`
    * **`repositories/`**: **(Lapisan Akses Data)** Mengenkapsulasi logika untuk operasi ke database (CRUD).
        * `budget_repository.go`: Akses data anggaran.
        * `expense_repository.go`: Akses data pengeluaran.
        * `user_repositories.go`: Akses data pengguna.
    * **`routes/`**: **(Lapisan Presentasi)** Mendefinisikan endpoint API dan menghubungkannya ke controller yang sesuai.
        * `routes.go`: Pengaturan semua rute aplikasi.
    * **`services/`**: **(Lapisan Logika Bisnis)** Mengimplementasikan aturan bisnis inti dan alur kerja aplikasi.
        * `auth_service.go`: Logika untuk autentikasi dan otorisasi.
        * `budget_service.go`: Operasi untuk manajemen anggaran.
        * `expense_service.go`: Operasi untuk pencatatan pengeluaran.
    * **`utils/`**: **(Lapisan Infrastruktur/Utilitas)** Berisi fungsi-fungsi pembantu yang dapat digunakan di berbagai lapisan.
        * `jwt_utils.go`: Utilitas untuk penanganan JSON Web Token (JWT).

* **`go.mod`**: Mendefinisikan modul Go dan dependensinya.
* **`go.sum`**: Berisi checksum dari dependensi yang digunakan.
* **`README.md`**: Dokumentasi utama proyek.
* **`tests/`**: Berisi kode pengujian untuk aplikasi.
    * `login_test.go`: Contoh file pengujian untuk fungsionalitas login.

---

### Manfaat Arsitektur Ini

* **Pemisahan Fungsi (Separation of Concerns)**: Setiap lapisan dan direktori memiliki tanggung jawab yang spesifik dan jelas.
* **Kemudahan Pemeliharaan (Maintainability)**: Perubahan pada satu lapisan atau modul cenderung tidak terlalu mempengaruhi bagian lain dari aplikasi.
* **Kemudahan Pengujian (Testability)**: Lebih mudah untuk menulis unit test untuk komponen yang terisolasi.
* **Alur Ketergantungan (Dependency Flow)**: Ketergantungan umumnya mengalir dari lapisan presentasi ke layanan, lalu ke repositori, menjaga logika inti tetap independen dari detail implementasi luar.
* **Antarmuka yang Bersih (Clean Interfaces)**: Setiap lapisan berkomunikasi melalui antarmuka yang terdefinisi dengan baik, mempromosikan kode yang lebih modular.
* **Enkapsulasi Kode Internal**: Penggunaan direktori `internal` mencegah penggunaan paket internal oleh proyek lain, menjaga integritas kode.

---

## Dokumentasi api
[Dokumentasi api](https://documenter.getpostman.com/view/39928139/2sB2qcBLVv)


## Struktur Proyek

```
.
├── cmd
│   └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   └── config.go
│   ├── controllers
│   │   ├── auth_controller.go
│   │   ├── budget_controller.go
│   │   └── expense_controller.go
│   ├── db
│   │   └── db.go
│   ├── dtos
│   │   ├── requests
│   │   │   ├── auth_request.go
│   │   │   ├── budget_request.go
│   │   │   └── expense_request.go
│   │   └── responses
│   │       ├── auth_response.go
│   │       ├── budget_response.go
│   │       ├── common_response.go
│   │       └── expense_response.go
│   ├── middlewares
│   │   └── auth_middleware.go
│   ├── models
│   │   ├── budget.go
│   │   ├── expense.go
│   │   └── user.go
│   ├── repositories
│   │   ├── budget_repository.go
│   │   ├── expense_repository.go
│   │   └── user_repositories.go
│   ├── routes
│   │   └── routes.go
│   ├── services
│   │   ├── auth_service.go
│   │   ├── budget_service.go
│   │   └── expense_service.go
│   └── utils
│       └── jwt_utils.go
├── README.md
└── tests
    └── login_test.go
```


## 📦 Instalasi

### Prasyarat

- Go 1.24 atau lebih baru
- MySQL 8.0 atau lebih baru
- Git

### Langkah Instalasi

1. Clone repository
```bash
git clone https://github.com/Alvarras/dompet-g0.git
cd dompet-g0
```

2. Install dependencies
```bash
go mod tidy
```

3. Setup database
```sql
CREATE DATABASE go_budget;
CREATE DATABASE go_budget_test;
```

4. Konfigurasi environment
```bash
# Copy file .env.example
cp .env.example .env

# Edit file .env sesuai konfigurasi
nano .env
```

### 🚀 Menjalankan Aplikasi
```bash
go run cmd/api/main.go
```
### Testing
```bash
go test -v ./...
```

