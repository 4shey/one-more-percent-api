# One More Percent 🚀

> Bot Telegram personal untuk tracking progress harian Ambatukam.  
> Reminder otomatis, deteksi completion via AI, ingatan percakapan, dan recap tengah malam.

---

## Tech Stack

| Layer | Teknologi | Versi |
|---|---|---|
| Language | Go | 1.26.3 |
| HTTP Server | `net/http` | stdlib |
| Database | PostgreSQL | 17.5-alpine |
| Telegram | `go-telegram-bot-api` | v5.5.1 |
| AI | Groq API (`llama-3.3-70b-versatile`) | — |
| DB Driver | `lib/pq` | v1.12.3 |
| Migration | `golang-migrate/migrate` | v4.18.3 |
| Hot Reload | `air` | v1.61.7 |
| Container | Docker + Docker Compose | — |

---

## Struktur Project

```
one_more_percent/
├── app/
│   └── main.go                          # Entry point: init DB → start scheduler → serve HTTP
│
├── internal/
│   ├── db/
│   │   └── db.go                        # Koneksi PostgreSQL via lib/pq
│   │
│   ├── models/
│   │   ├── schedule.go                  # Struct Schedule
│   │   └── progress.go                  # Struct ScheduleProgress (dengan joined fields)
│   │
│   ├── services/
│   │   ├── ai_service.go                # callGroq, AskAI, DetectCompletion, GenerateReminder,
│   │   │                                #   GenerateRecap, GenerateCompletionReply
│   │   ├── conversation_service.go      # In-memory rolling history 10 pesan per chatID
│   │   ├── telegram_service.go          # SendTelegramMessage
│   │   ├── progress_service.go          # CRUD schedule_progresses, GetActiveSchedule,
│   │   │                                #   ProgressRowExists
│   │   └── scheduler_service.go         # Ticker per menit, catch-up check, midnight recap,
│   │                                    #   in-memory pendingReminders map
│   │
│   ├── handlers/
│   │   ├── telegram_handler.go          # Webhook handler: pending check → intent detect / chat
│   │   └── healt_handler.go             # GET / health check
│   │
│   └── routes/
│       └── routes.go                    # Route registration
│
├── database/
│   ├── migrations/
│   │   ├── 000001_create_users.up.sql
│   │   ├── 000001_create_users.down.sql
│   │   ├── 000002_create_schedules.up.sql
│   │   ├── 000002_create_schedules.down.sql
│   │   ├── 000003_create_schedule_progresses.up.sql
│   │   └── 000003_create_schedule_progresses.down.sql
│   │
│   └── seeders/
│       └── seed_schedules.sql           # User Ambatukam + jadwal Senin
│
├── .air.toml                            # Konfigurasi hot reload
├── .env                                 # Environment variables (tidak di-commit)
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
```

---

## Database Schema

### `users`
```sql
CREATE TABLE users (
    id          SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE,
    name        TEXT NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW()
);
```

### `schedules`
```sql
CREATE TABLE schedules (
    id          SERIAL PRIMARY KEY,
    user_id     INT NOT NULL REFERENCES users(id),
    day_of_week VARCHAR(10) NOT NULL,   -- "Monday" | "Tuesday" | dst
    start_time  TIME NOT NULL,
    end_time    TIME NOT NULL,
    activity    VARCHAR(100) NOT NULL,
    is_active   BOOLEAN DEFAULT TRUE,
    created_at  TIMESTAMP DEFAULT NOW()
);
```

### `schedule_progresses`
```sql
CREATE TABLE schedule_progresses (
    id            SERIAL PRIMARY KEY,
    schedule_id   INT NOT NULL REFERENCES schedules(id),
    progress_date DATE NOT NULL,
    status        VARCHAR(20) NOT NULL DEFAULT 'pending',
    --            ↑ pending | completed | missed
    completed_at  TIMESTAMP NULL,
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW(),

    UNIQUE(schedule_id, progress_date)   -- 1 jadwal = 1 row per hari
);
```

> **Reset system**: tidak perlu truncate. Besok = progress_date baru = row baru secara implisit. History aman untuk analytics.

---

## Alur Sistem Lengkap

### A. Startup Catch-Up

Saat server pertama kali start (atau restart):

```
StartScheduler()
    │
    ├─ LoadLocation("Asia/Jakarta")
    │
    └─ runCatchUpCheck()
           │
           ├─ GetActiveSchedule(dayName, HH:MM)
           │       │
           │       ├─ Tidak ada jadwal aktif sekarang?
           │       │     → log "nothing to catch up", selesai
           │       │
           │       └─ Ada jadwal aktif (misal: Tidur 01:00-05:00)?
           │             │
           │             ├─ ProgressRowExists(scheduleID, today)?
           │             │     → YES: sudah diremind, skip
           │             │     → NO:  sendReminder() ← kirim sekarang
           │
           └─ [sleep hingga menit berikutnya]
```

**Tujuan**: Jika server restart saat sedang dalam window jadwal, reminder tetap terkirim tanpa harus menunggu menit berikutnya.

---

### B. Reminder Rutin (setiap menit)

```
[Ticker setiap 60 detik, aligned ke batas menit]
    │
    └─ runCheck()
           │
           ├─ Jam 00:00? → runMidnightRecap(kemarin)
           │
           └─ GetSchedulesForDay(dayName)
                   │
                   └─ Loop semua jadwal
                           │
                           └─ StartTime == HH:MM sekarang?
                                   │
                                   └─ sendReminder(schedule, now)
                                           │
                                           ├─ EnsureProgressRow (status=pending)
                                           ├─ SetPendingReminder(chatID, scheduleInfo)
                                           └─ GenerateReminder → SendTelegramMessage
```

---

### C. User Reply → Completion Detection

```
[Telegram kirim update ke /webhook]
    │
    └─ TelegramWebhookHandler
           │
           ├─ GetPendingReminder(chatID)
           │
           ├─ Ada pending reminder?
           │     │
           │     YES
           │     │
           │     └─ DetectCompletion(userMessage)   ← AI intent classification
           │             │
           │             ├─ YES ("udah", "done", "selesai", dll)
           │             │       ├─ MarkCompleted(scheduleID, date)
           │             │       ├─ ClearPendingReminder(chatID)
           │             │       └─ GenerateCompletionReply → kirim ke Telegram
           │             │
           │             └─ NO ("belum", "nanti", dll)
           │                     └─ AskAI(chatID, message) → normal chat
           │
           └─ Tidak ada pending?
                   └─ AskAI(chatID, message) → normal chat
```

---

### D. AskAI — Percakapan dengan Context

Setiap kali `AskAI` dipanggil:

```
AskAI(chatID, message)
    │
    ├─ getActiveScheduleContext()
    │       └─ Query: ada jadwal start ≤ now < end hari ini?
    │               ├─ Ada  → "[Konteks: Jadwal aktif 'Golang' (13:00-15:00)]"
    │               └─ Tidak → string kosong
    │
    ├─ buildSystemPrompt(scheduleCtx)
    │       └─ Profil user (Ambatukam, mahasiswa, remote) + konteks jadwal
    │
    ├─ GetHistory(chatID)
    │       └─ Slice 10 pesan terakhir (rolling, thread-safe)
    │
    ├─ callGroq([system, ...history, user_message])
    │
    ├─ AddMessage(chatID, "user", message)
    └─ AddMessage(chatID, "assistant", reply)
```

**Contoh messages yang dikirim ke Groq:**
```
[system]    Kamu adalah One More Percent, asisten Ambatukam.
            Ambatukam adalah mahasiswa yang mengejar kerja remote.
            [Konteks: Sekarang jadwal aktif adalah 'Tidur' (01:00-05:00).]

[user]      halo bang
[assistant] halo! lagi tidur atau belum?
[user]      belum, masih scrolling
[assistant] ya ampun bang, tidur dulu sana 😴
[user]      iya deh          ← pesan sekarang
```

---

### E. Midnight Recap (00:00 WIB)

```
runMidnightRecap(yesterday)
    │
    ├─ MarkAllPendingAsMissed(yesterday)
    │       └─ UPDATE status='missed' WHERE status='pending' AND date=yesterday
    │
    ├─ GetDayProgress(yesterday)
    │       └─ JOIN schedules + schedule_progresses
    │               → completed: [English, Golang]
    │               → missed:    [Portfolio]
    │
    └─ GenerateRecap(completed, missed)
            └─ AI generate recap santai → SendTelegramMessage
```

**Contoh output recap:**
```
hari ini lumayan bang 😼

✅ tidur
✅ english
✅ golang
❌ portfolio

3 dari 4 selesai.
ga perfect gapapa, besok lanjut lagi
one more percent 🚀
```

---

## AI Functions

| Fungsi | Input | Output | Dipakai di |
|---|---|---|---|
| `callGroq([]Message)` | Full message list | Raw AI string | Semua fungsi di bawah |
| `AskAI(chatID, msg)` | chatID + pesan user | Balasan chat | Webhook handler |
| `DetectCompletion(msg)` | Pesan user | `true` / `false` | Webhook handler |
| `GenerateReminder(activity)` | Nama aktivitas | Teks reminder | Scheduler |
| `GenerateRecap(completed, missed)` | Slice progress | Teks recap | Midnight scheduler |
| `GenerateCompletionReply(activity)` | Nama aktivitas | Teks konfirmasi | Webhook handler |

### DetectCompletion — Contoh

```
"udah bang"            → YES
"baru selesai english" → YES
"done"                 → YES
"finished"             → YES

"belum"                → NO
"nanti dulu"           → NO
"masih males"          → NO
"ga jadi"              → NO
```

---

## Conversation Memory

- Disimpan **in-memory** per `chatID` (map + mutex)
- Rolling window: **10 pesan terakhir** (5 user + 5 assistant)
- Diisi setiap call `AskAI`
- Reset saat server restart (by design — tidak perlu persistent untuk bot personal)

```
File: internal/services/conversation_service.go

AddMessage(chatID, role, content)  → tambah ke history
GetHistory(chatID) []Message       → ambil salinan history
```

---

## Environment Variables (`.env`)

```env
# Telegram
BOT_TOKEN=<telegram_bot_token>
TELEGRAM_CHAT_ID=<your_telegram_chat_id>

# Database
DB_HOST=db
DB_PORT=5432
DB_USER=omp
DB_PASSWORD=omp123
DB_NAME=onemorepercent

# AI — Groq
MODEL_TOKEN=<groq_api_key>
MODEL_NAME=llama-3.3-70b-versatile
```

---

## Setup & Menjalankan

### Prasyarat

- Docker & Docker Compose
- Telegram Bot Token → [BotFather](https://t.me/BotFather)
- Groq API Key → [console.groq.com](https://console.groq.com)

### Langkah-langkah

```bash
# 1. Isi environment variables
cp .env.example .env
# edit .env: isi BOT_TOKEN, TELEGRAM_CHAT_ID, MODEL_TOKEN

# 2. Jalankan database
docker compose up -d db

# 3. Jalankan migrasi (buat semua tabel)
docker compose --profile tools run --rm migrate

# 4. Seed data awal (user + jadwal)
docker compose exec -T db psql -U omp -d onemorepercent \
  < database/seeders/seed_schedules.sql

# 5. Jalankan aplikasi
docker compose up -d app

# 6. Cek logs
docker compose logs -f app
```

### Set Telegram Webhook

```bash
# Production (pakai domain public)
curl -X POST "https://api.telegram.org/bot<BOT_TOKEN>/setWebhook" \
  -d "url=https://<your-domain>/webhook"

# Development (pakai ngrok)
ngrok http 8080
# lalu set webhook ke URL ngrok yang muncul
```

---

## Mengelola Jadwal

### Tambah jadwal baru

```bash
docker compose exec -T db psql -U omp -d onemorepercent -c "
  INSERT INTO schedules (user_id, day_of_week, start_time, end_time, activity)
  VALUES (1, 'Tuesday', '09:00', '10:00', 'Reading');
"
```

Jadwal langsung aktif — scheduler baca DB setiap tick, tidak perlu restart.

### Nonaktifkan jadwal

```bash
docker compose exec -T db psql -U omp -d onemorepercent -c "
  UPDATE schedules SET is_active = FALSE WHERE id = 3;
"
```

### Lihat semua jadwal

```bash
docker compose exec -T db psql -U omp -d onemorepercent -c "
  SELECT id, day_of_week,
         TO_CHAR(start_time,'HH24:MI') AS start,
         TO_CHAR(end_time,'HH24:MI') AS end,
         activity, is_active
  FROM schedules ORDER BY day_of_week, start_time;
"
```

### Nilai `day_of_week` yang valid

```
Monday | Tuesday | Wednesday | Thursday | Friday | Saturday | Sunday
```

---

## HTTP Endpoints

| Method | Path | Keterangan |
|---|---|---|
| `GET` | `/` | Health check — returns 200 OK |
| `POST` | `/webhook` | Telegram webhook receiver |

---

## Development

Project pakai **Air** untuk hot reload — file berubah → otomatis rebuild & restart.

```bash
# Log real-time
docker compose logs -f app

# Rebuild manual
docker compose up -d --build app

# Reset database total (hapus semua data + volume)
docker compose down -v

# Masuk ke psql
docker compose exec db psql -U omp -d onemorepercent
```

---

## Dependencies

```
github.com/go-telegram-bot-api/telegram-bot-api/v5  v5.5.1
github.com/lib/pq                                   v1.12.3
```

Migration tool (Docker only, profile `tools`):
```
migrate/migrate  v4.18.3
```
