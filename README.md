## InvGate Armmada Portal

Portal end-user berbasis Vue 3 + Go untuk membuat dan memantau ticket InvGate Armmada.

### Struktur Proyek

```
backend/   # Go clean architecture (auth, ticket, invgate client)
frontend/  # Vue 3 + TypeScript + Carbon Design System
docker-compose.yml
```

Backend (Gin + Gorm + MySQL):

- `cmd/server`: entrypoint HTTP server.
- `internal/config`: konfigurasi environment.
- `internal/auth`, `internal/user`, `internal/ticket`: domain logic & handler (registrasi user otomatis disinkronkan ke InvGate dan DB lokal).
- `internal/invgate`: client HTTP ke API InvGate Armmada.
- `internal/middleware`: logging, error recovery, JWT auth.
- `migrations`: migrasi MySQL (tabel `users`).

Frontend:

- Vue 3 + Vite + TypeScript.
- Folder `src/api`, `src/composables`, `src/components`, `src/pages`, `src/router`, `src/stores`, `src/utils`.
- Carbon Design System untuk UI dan TanStack Query untuk fetching + polling.

### Konfigurasi Environment

1. Salin `backend/env.example` menjadi `backend/.env` dan sesuaikan.
2. Salin `frontend/env.example` menjadi `frontend/.env` bila perlu override `VITE_API_BASE_URL` atau menonaktifkan dummy data (`VITE_USE_DUMMY=false` untuk koneksi backend nyata).

Variabel penting:

```
SERVER_PORT
DB_HOST / DB_PORT / DB_USER / DB_PASSWORD / DB_NAME
JWT_SECRET
ARMMADA_BASE_URL / ARMMADA_USERNAME / ARMMADA_PASSWORD / ARMMADA_PAGE_KEY
VITE_API_BASE_URL / VITE_USE_DUMMY
```

### Menjalankan Secara Lokal

```bash
cd backend
go run main.go

cd ../frontend
npm install
npm run dev
```

Backend berjalan di `http://localhost:8080`, frontend dev di `http://localhost:5173`.

### Jalankan via Docker Compose

```bash
docker compose up --build
```

Layanan:

- `frontend`: http://localhost:5173
- `backend`: http://localhost:8080/api
- `db`: MySQL 8 (port 3306)

### Endpoint Backend

| Method | Path            | Deskripsi                  |
|--------|-----------------|---------------------------|
| POST   | `/api/auth/register` | Registrasi user baru (sinkron ke InvGate) |
| POST   | `/api/auth/login`    | Login (JWT)              |
| POST   | `/api/tickets`       | Kirim ticket ke InvGate  |
| GET    | `/api/tickets`       | Daftar ticket user       |
| GET    | `/api/tickets/{id}`  | Detail ticket + komentar |

Semua endpoint `/api/tickets` membutuhkan header `Authorization: Bearer <token>`.

### TanStack Query Interval

- Daftar ticket (`useTicketList`) otomatis refresh setiap 30 detik.
- Detail ticket (`useTicketDetail`) refresh setiap 15 detik untuk tracking status.

### Contoh Pemanggilan InvGate (Backend)

```go
payload := invgate.CreateTicketPayload{
    SourceID:   1,
    CreatorID:  100,
    CustomerID: 100,
    CategoryID: 10,
    TypeID:     20,
    PriorityID: 30,
    Title:      "Laptop bermasalah",
    Description:"Keyboard tidak berfungsi",
}
resp, err := invgateService.CreateTicket(ctx, payload)
```

Service menggunakan `net/http` + Basic Auth (`ARMMADA_USERNAME/PASSWORD`) dan otomatis menambahkan `ARMMADA_PAGE_KEY` untuk pagination saat mengambil daftar ticket.

