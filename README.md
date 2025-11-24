## 🚀 Arsitektur

Service untuk manajemen nomor telepon dengan pendekatan **Hexagonal Architecture (Ports & Adapters)**.  
Tujuannya adalah menjaga **domain/business logic tetap murni**, terpisah dari framework, database, dan teknologi eksternal lainnya.

Dependency dalam arsitektur hexagonal mengikuti Dependency Inversion Principle:
- Logika domain bergantung pada abstraksi (ports), bukan implementasi konkret (adapters).
- Logika domain dan adapters berada di paket/module terpisah.
- Domain berada di pusat dan tidak mengetahui keberadaan adapters.
- Adapters bergantung pada ports yang mereka implementasikan, bukan sebaliknya.

## 🚀 Menjalankan Project

### 1. Clone Repository
```bash
  git clone https://github.com/username/ms-telnum-manager.git
  cd clean-architecture
```

### 2. Install Dependency
```bash
  go mod tidy
```

### 3. Konfigurasi Env
```bash
  .env
```

### 4. Jalankan Semua Migrasi
```bash
  go run main.go migrate up
```

### 5. Rollback Migrasi
```bash
  go run main.go migrate down
```

### 6. Cek Migrasi Status
```bash
  go run main.go migrate status
```

### 7. Menjalankan Service
```bash
  go run main.go start
```

---

### 🔄 Gambaran Arsitektur
                [ Inbound Adapters ]
         ┌───────────────────────────────────┐
         │   HTTP (Echo)   |   WebSocket     │
         │ ← internal/adapter/inbound/...    │
         │     request/response DTO          │
         └───────────────────────────────────┘
                          │
                          ▼
               ┌─────────────────────────┐
               │       Application       │
               │   (app/ - dependency    │
               │    injection, orches.)  │
               └─────────────────────────┘
                          │
               ┌─────────────────────────┐
               │          Ports          │
               │ (internal/port/... →    │
               │  inbound & outbound)    │
               └─────────────────────────┘
                          │
               ┌─────────────────────────┐
               │         Domain          │
               │ entity/  +  service/    │
               │ (logika bisnis murni)   │
               └─────────────────────────┘
                          │
                          ▼
                [ Outbound Adapters ]
     ┌────────────────┬────────────────────────┐
     │ Repository (DB)│  HTTP Client / Broker  │
     │ postgres/...   │ httpclient/...         │
     │ model/,repo/   │ (Kafka, Rabbit, dll)   │
     └────────────────┴────────────────────────┘

---

## 📂 Struktur Project

```bash
ms-telnum-manager/
├── cmd/                           # Entry point (main.go, bootstrap)
│
├── config/                        # File konfigurasi (env, yaml, json)
│
├── internal/                      # Semua kode inti aplikasi
│   ├── adapter/                   # Implementasi inbound & outbound adapter
│   │   ├── inbound/               # Apa yang aplikasi tawarkan (HTTP, gRPC, WebSocket, dll.)
│   │   │   ├── echo/              
│   │   │   │   ├── request/       # DTO request untuk Echo
│   │   │   │   ├── response/      # DTO response untuk Echo
│   │   │
│   │   ├── outbound/              # Apa yang aplikasi butuhkan (DB, broker, storage, API eksternal)
│   │       ├── httpclient/        # Adapter untuk komunikasi HTTP dengan service eksternal
│   │       ├── postgres/          
│   │           ├── model/         # Struktur data mapping tabel database
│   │           ├── repository/    # Implementasi repository (akses DB)
│   │
│   ├── app/                       # Dependency injection (menyambungkan domain, ports, adapters)
│   │
│   ├── domain/                    # Bisnis logic inti (bersih, tanpa ketergantungan eksternal)
│   │   ├── entity/                # Definisi entity utama (misal: User, TelNumber)
│   │   ├── service/               # Aturan & proses bisnis (use case)
│   │
│   ├── migration/                 # File migrasi database (create/drop table)
│   │
│   ├── port/                      # Interface (kontrak komunikasi domain ↔ adapter)
│       ├── inbound/               # Interface inbound adapter (contoh: handler contract)
│       ├── outbound/              # Interface outbound adapter (contoh: repository contract)
│
├── tests/                         # Unit test & integration test
│
├── utils/                         # Kumpulan helper umum yang reusable
│   ├── conv/                      # Utility konversi tipe data
│   ├── ping/                      # Utility health check
│   ├── validator/                 # Utility validasi input
│   ├── encryption/                # Utility enkripsi/dekripsi
│
├── .env                           # Variabel environment
├── .gitignore                     # Ignore file untuk Git
├── go.mod                         # Go module dependencies
├── main.go                        # Entry point utama

```
