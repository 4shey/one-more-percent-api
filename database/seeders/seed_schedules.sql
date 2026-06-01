-- =============================================================
-- Seed: One More Percent
-- User: Ambatukam | Mahasiswa yang mengejar kerja remote
-- =============================================================

-- 1. User
INSERT INTO users (telegram_id, name)
VALUES (6616220735, 'Ambatukam')
ON CONFLICT (telegram_id) DO UPDATE SET name = EXCLUDED.name;

-- 2. Jadwal Senin (Monday)
INSERT INTO schedules (user_id, day_of_week, start_time, end_time, activity)
VALUES
  (1, 'Monday', '01:00', '05:00', 'Tidur'),
  (1, 'Monday', '11:00', '12:00', 'English'),
  (1, 'Monday', '13:00', '15:00', 'Golang'),
  (1, 'Monday', '20:00', '21:00', 'Portfolio')
ON CONFLICT DO NOTHING;

-- Tambahkan jadwal hari lain di sini:
-- INSERT INTO schedules (user_id, day_of_week, start_time, end_time, activity)
-- VALUES
--   (1, 'Tuesday', '09:00', '10:00', 'Reading'),
--   (1, 'Wednesday', '13:00', '15:00', 'Golang');
