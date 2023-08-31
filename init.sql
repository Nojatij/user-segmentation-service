-- init.sql

-- Создание таблицы сегментов
CREATE TABLE segments (
    id SERIAL PRIMARY KEY,
    slug TEXT NOT NULL
);

-- Вставка тестовых данных
INSERT INTO segments (slug) VALUES
    ('AVITO_VOICE_MESSAGES'),
    ('AVITO_PERFORMANCE_VAS'),
    ('AVITO_DISCOUNT_30'),
    ('AVITO_DISCOUNT_50');
