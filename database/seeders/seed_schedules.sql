-- =============================================================
-- SEED: ONE MORE PERCENT (FIXED VERSION)
-- =============================================================

WITH u AS (
    INSERT INTO users (telegram_id, name)
    VALUES (6616220735, 'Ambatukam')
    ON CONFLICT (telegram_id)
    DO UPDATE SET name = EXCLUDED.name
    RETURNING id
)

INSERT INTO schedules (user_id, day_of_week, start_time, end_time, activity)
SELECT u.id, v.day_of_week, v.start_time::time, v.end_time::time, v.activity
FROM u
CROSS JOIN (VALUES

-- =============================================================
-- MONDAY
-- =============================================================
('Monday', '04:30', '06:00', 'Olahraga'),
('Monday', '07:00', '11:30', 'Kuliah / Belajar Produktif'),
('Monday', '13:00', '16:00', 'Golang & Backend'),
('Monday', '16:00', '17:00', 'Olahraga'),
('Monday', '17:00', '20:00', 'English'),
('Monday', '20:00', '21:00', 'Cyber Security'),
('Monday', '21:00', '22:00', 'Web Development'),
('Monday', '22:00', '23:00', 'AI Engineer'),

-- =============================================================
-- TUESDAY
-- =============================================================
('Tuesday', '04:30', '06:00', 'Olahraga'),
('Tuesday', '07:00', '11:30', 'Kuliah / Belajar Produktif'),
('Tuesday', '13:00', '16:00', 'Web Development'),
('Tuesday', '16:00', '17:00', 'Olahraga'),
('Tuesday', '19:00', '20:00', 'English'),
('Tuesday', '20:00', '21:00', 'Data Science'),
('Tuesday', '21:00', '22:00', 'Cyber Security'),
('Tuesday', '22:00', '23:00', 'AI Engineer'),

-- =============================================================
-- WEDNESDAY
-- =============================================================
('Wednesday', '04:30', '06:00', 'Olahraga'),
('Wednesday', '07:00', '11:30', 'Kuliah / Belajar Produktif'),
('Wednesday', '13:00', '16:00', 'Data Science'),
('Wednesday', '16:00', '17:00', 'Olahraga'),
('Wednesday', '19:00', '20:00', 'English'),
('Wednesday', '20:00', '21:00', 'Web Development'),
('Wednesday', '21:00', '22:00', 'Cyber Security'),
('Wednesday', '22:00', '23:00', 'AI Engineer'),

-- =============================================================
-- THURSDAY
-- =============================================================
('Thursday', '04:30', '06:00', 'Olahraga'),
('Thursday', '07:00', '11:30', 'Kuliah / Belajar Produktif'),
('Thursday', '13:00', '16:00', 'AI Engineer'),
('Thursday', '16:00', '17:00', 'Olahraga'),
('Thursday', '19:00', '20:00', 'English'),
('Thursday', '20:00', '21:00', 'Cyber Security'),
('Thursday', '21:00', '22:00', 'Web Development'),
('Thursday', '22:00', '23:00', 'Data Science'),

-- =============================================================
-- FRIDAY
-- =============================================================
('Friday', '04:30', '06:00', 'Olahraga'),
('Friday', '07:00', '11:30', 'Kuliah / Belajar Produktif'),
('Friday', '13:00', '16:00', 'Cyber Security'),
('Friday', '16:00', '17:00', 'Olahraga'),
('Friday', '19:00', '20:00', 'English'),
('Friday', '20:00', '21:00', 'AI Engineer'),
('Friday', '21:00', '22:00', 'Web Development'),
('Friday', '22:00', '23:00', 'Data Science'),

-- =============================================================
-- SATURDAY
-- =============================================================
('Saturday', '04:30', '06:00', 'Olahraga'),
('Saturday', '07:00', '11:30', 'Project / Deep Work'),
('Saturday', '13:00', '16:00', 'Portfolio & Project'),
('Saturday', '16:00', '17:00', 'Olahraga'),
('Saturday', '19:00', '20:00', 'English'),
('Saturday', '20:00', '21:00', 'Cyber Security'),
('Saturday', '21:00', '22:00', 'AI Engineer'),
('Saturday', '22:00', '23:00', 'Web Development'),

-- =============================================================
-- SUNDAY
-- =============================================================
('Sunday', '04:30', '06:00', 'Olahraga'),
('Sunday', '07:00', '11:30', 'Review Mingguan / Belajar'),
('Sunday', '13:00', '16:00', 'Portfolio & Planning'),
('Sunday', '16:00', '17:00', 'Olahraga'),
('Sunday', '19:00', '20:00', 'English'),
('Sunday', '20:00', '21:00', 'Data Science'),
('Sunday', '21:00', '22:00', 'Cyber Security'),
('Sunday', '22:00', '23:00', 'AI Engineer')

) AS v(day_of_week, start_time, end_time, activity);