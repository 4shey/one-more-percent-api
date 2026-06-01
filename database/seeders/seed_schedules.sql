-- =============================================================
-- Seed: One More Percent
-- User: Ambatukam | Mahasiswa yang mengejar kerja remote
-- =============================================================

-- 1. User
INSERT INTO users (telegram_id, name) VALUES (6616220735, 'Ambatukam') ON CONFLICT (telegram_id) DO UPDATE SET name = EXCLUDED.name;

-- 2. Jadwal Senin (Monday)
INSERT INTO schedules ( user_id, day_of_week, start_time, end_time, activity ) VALUES
-- =============================================================
-- Monday
-- =============================================================
(1, 'Monday', '04:30', '06:00', 'Olahraga'),
(1, 'Monday', '07:00', '11:30', 'Kuliah / Belajar Produktif'),
(1, 'Monday', '13:00', '16:00', 'Golang & Backend'),
(1, 'Monday', '16:00', '17:00', 'Olahraga'),
(1, 'Monday', '19:00', '20:00', 'English'),
(1, 'Monday', '20:00', '21:00', 'Cyber Security'),
(1, 'Monday', '21:00', '22:00', 'Web Development'),
(1, 'Monday', '22:00', '23:00', 'AI Engineer'),

-- =============================================================
-- Tuesday
-- =============================================================
(1, 'Tuesday', '04:30', '06:00', 'Olahraga'),
(1, 'Tuesday', '07:00', '11:30', 'Kuliah / Belajar Produktif'),
(1, 'Tuesday', '13:00', '16:00', 'Web Development'),
(1, 'Tuesday', '16:00', '17:00', 'Olahraga'),
(1, 'Tuesday', '19:00', '20:00', 'English'),
(1, 'Tuesday', '20:00', '21:00', 'Data Science'),
(1, 'Tuesday', '21:00', '22:00', 'Cyber Security'),
(1, 'Tuesday', '22:00', '23:00', 'AI Engineer'),

-- =============================================================
-- Wednesday
-- =============================================================
(1, 'Wednesday', '04:30', '06:00', 'Olahraga'),
(1, 'Wednesday', '07:00', '11:30', 'Kuliah / Belajar Produktif'),
(1, 'Wednesday', '13:00', '16:00', 'Data Science'),
(1, 'Wednesday', '16:00', '17:00', 'Olahraga'),
(1, 'Wednesday', '19:00', '20:00', 'English'),
(1, 'Wednesday', '20:00', '21:00', 'Web Development'),
(1, 'Wednesday', '21:00', '22:00', 'Cyber Security'),
(1, 'Wednesday', '22:00', '23:00', 'AI Engineer'),

-- =============================================================
-- Thursday
-- =============================================================
(1, 'Thursday', '04:30', '06:00', 'Olahraga'),
(1, 'Thursday', '07:00', '11:30', 'Kuliah / Belajar Produktif'),
(1, 'Thursday', '13:00', '16:00', 'AI Engineer'),
(1, 'Thursday', '16:00', '17:00', 'Olahraga'),
(1, 'Thursday', '19:00', '20:00', 'English'),
(1, 'Thursday', '20:00', '21:00', 'Cyber Security'),
(1, 'Thursday', '21:00', '22:00', 'Web Development'),
(1, 'Thursday', '22:00', '23:00', 'Data Science'),

-- =============================================================
-- Friday
-- =============================================================
(1, 'Friday', '04:30', '06:00', 'Olahraga'),
(1, 'Friday', '07:00', '11:30', 'Kuliah / Belajar Produktif'),
(1, 'Friday', '13:00', '16:00', 'Cyber Security'),
(1, 'Friday', '16:00', '17:00', 'Olahraga'),
(1, 'Friday', '19:00', '20:00', 'English'),
(1, 'Friday', '20:00', '21:00', 'AI Engineer'),
(1, 'Friday', '21:00', '22:00', 'Web Development'),
(1, 'Friday', '22:00', '23:00', 'Data Science'),

-- =============================================================
-- Saturday
-- =============================================================
(1, 'Saturday', '04:30', '06:00', 'Olahraga'),
(1, 'Saturday', '07:00', '11:30', 'Project / Deep Work'),
(1, 'Saturday', '13:00', '16:00', 'Portfolio & Project'),
(1, 'Saturday', '16:00', '17:00', 'Olahraga'),
(1, 'Saturday', '19:00', '20:00', 'English'),
(1, 'Saturday', '20:00', '21:00', 'Cyber Security'),
(1, 'Saturday', '21:00', '22:00', 'AI Engineer'),
(1, 'Saturday', '22:00', '23:00', 'Web Development'),

-- =============================================================
-- Sunday
-- =============================================================
(1, 'Sunday', '04:30', '06:00', 'Olahraga'),
(1, 'Sunday', '07:00', '11:30', 'Review Mingguan / Belajar'),
(1, 'Sunday', '13:00', '16:00', 'Portfolio & Planning'),
(1, 'Sunday', '16:00', '17:00', 'Olahraga'),
(1, 'Sunday', '19:00', '20:00', 'English'),
(1, 'Sunday', '20:00', '21:00', 'Data Science'),
(1, 'Sunday', '21:00', '22:00', 'Cyber Security'),
(1, 'Sunday', '22:00', '23:00', 'AI Engineer')

ON CONFLICT DO NOTHING;
```
