# One More Percent

Personal Telegram bot berbasis AI untuk membantu memantau progres harian melalui pengingat otomatis dan percakapan natural.

Bot dapat mengingatkan jadwal, memahami respon, lalu mencatat progres secara otomatis tanpa perintah khusus.

---

## Table of Contents

* [Fitur](#fitur)
* [Teknologi](#teknologi)
* [Menjalankan Project](#menjalankan-project)
* [Mengelola Jadwal](#mengelola-jadwal)
* [Contoh Interaksi](#contoh-interaksi)

---

## Fitur

* Pengingat otomatis berdasarkan jadwal
* Deteksi penyelesaian tugas melalui percakapan biasa
* AI memahami konteks jadwal aktif
* Rekap aktivitas harian setiap pukul `00:00 WIB`

---

## Teknologi

| Component      | Stack                                |
| -------------- | ------------------------------------ |
| Backend        | Go 1.26                              |
| Database       | PostgreSQL 17.5                      |
| AI             | Groq API (`llama-3.3-70b-versatile`) |
| Messaging      | Telegram Bot API                     |
| Infrastructure | Docker & Docker Compose              |
| Migration      | Golang Migrate                       |

---

## Menjalankan Project

### 1. Salin file environment

```bash
cp .env.example .env
```

Setelah itu, sesuaikan konfigurasi pada file `.env`.

Konfigurasi yang perlu diperhatikan:

* Telegram Bot Token
* Telegram User ID
* Groq API Key
* Database configuration

### 2. Jalankan database

```bash
docker compose up -d db
```

### 3. Jalankan migration

```bash
docker compose --profile tools run --rm migrate
```

Migration digunakan untuk membuat struktur tabel pada database.

### 4. Tambahkan data awal (opsional)

```bash
docker compose exec -T db psql -U omp -d onemorepercent < database/seeders/seed_schedules.sql
```

Langkah ini bersifat opsional jika ingin langsung menggunakan data contoh.

### 5. Jalankan aplikasi

```bash
docker compose up -d app
```

Untuk melihat log aplikasi:

```bash
docker compose logs -f app
```

### 6. Set webhook Telegram

Pastikan aplikasi dapat diakses dari internet menggunakan domain, reverse proxy, atau tunnel.

Kemudian jalankan:

```bash
curl -X POST "https://api.telegram.org/bot<BOT_TOKEN>/setWebhook" \
-d "url=https://your-domain.com/webhook"
```

Contoh:

```bash
curl -X POST "https://api.telegram.org/bot123456:ABCDEF/setWebhook" \
-d "url=https://example.com/webhook"
```

---

## Mengelola Jadwal

Saat ini jadwal dikelola langsung melalui database.

Perubahan jadwal akan terbaca otomatis tanpa perlu me-restart aplikasi.

### Menambahkan jadwal

```sql
INSERT INTO schedules (
    day_of_week,
    start_time,
    end_time,
    activity
)
VALUES (
    'Tuesday',
    '09:00',
    '10:00',
    'Membaca Buku'
);
```

Atau melalui terminal:

```bash
docker compose exec -T db psql -U omp -d onemorepercent -c "
INSERT INTO schedules (
  day_of_week,
  start_time,
  end_time,
  activity
)
VALUES (
  'Tuesday',
  '09:00',
  '10:00',
  'Membaca Buku'
);
"
```

### Format hari

Gunakan format bahasa Inggris berikut:

| Indonesia | Format Database |
| --------- | --------------: |
| Senin     |          Monday |
| Selasa    |         Tuesday |
| Rabu      |       Wednesday |
| Kamis     |        Thursday |
| Jumat     |          Friday |
| Sabtu     |        Saturday |
| Minggu    |          Sunday |

---