# One More Percent

Bot Telegram personal berbasis AI untuk melacak progres harian kamu. Bot ini otomatis mengingatkan jadwalmu dan mendeteksi apakah suatu tugas sudah selesai atau belum melalui obrolan biasa (natural language).

## Fitur Utama
- **Auto-Reminder**: Notifikasi Telegram otomatis sesuai jadwal yang kamu tentukan.
- **AI Completion**: Bot paham bahasa manusia. Cukup bilang "udah kelar" atau "baru selesai", jadwalmu akan otomatis ditandai selesai.
- **Context-Aware AI**: Asisten AI tahu jadwalmu saat ini dan merespons obrolan sesuai konteks dengan persona yang fokus pada produktivitas.
- **Daily Recap**: Setiap tengah malam (00:00 WIB), bot mengirim rangkuman jadwal yang berhasil dan gagal dikerjakan.

## Tech Stack
- **Go 1.26** & `net/http`
- **PostgreSQL 17.5** (`lib/pq` & `golang-migrate`)
- **Telegram Bot API** (`go-telegram-bot-api/v5`)
- **Groq API** (`llama-3.3-70b-versatile`) untuk asisten AI
- **Docker & Docker Compose**

## Cara Menjalankan (via Docker)

1. **Persiapan Environment**
   Copy file `.env.example` dan isi konfigurasinya (terutama Token Bot, ID Telegram, dan API Key Groq).
   ```bash
   cp .env.example .env
   ```

2. **Jalankan Database & Migrasi**
   ```bash
   docker compose up -d db
   docker compose --profile tools run --rm migrate
   ```

3. **Isi Data Awal (Opsional)**
   ```bash
   docker compose exec -T db psql -U omp -d onemorepercent < database/seeders/seed_schedules.sql
   ```

4. **Jalankan Aplikasi**
   ```bash
   docker compose up -d app
   ```

5. **Set Webhook Telegram**
   Pastikan port 8080 terhubung ke internet (bisa menggunakan domain atau ngrok) lalu set webhook:
   ```bash
   curl -X POST "https://api.telegram.org/bot<BOT_TOKEN>/setWebhook" -d "url=https://<your-domain>/webhook"
   ```

## Mengelola Jadwal

Saat ini manajemen jadwal dilakukan langsung di database. Bot akan otomatis mendeteksi perubahan tanpa perlu direstart.

**Contoh tambah jadwal:**
```bash
docker compose exec -T db psql -U omp -d onemorepercent -c "
  INSERT INTO schedules (day_of_week, start_time, end_time, activity)
  VALUES ('Tuesday', '09:00', '10:00', 'Membaca Buku');
"
```
*(Catatan: Hari menggunakan format bahasa Inggris: Monday, Tuesday, dst.)*
