-- Создание базы данных
CREATE DATABASE my_time_tracker;

-- Подключение к базе данных
\c my_time_tracker;

-- Создание таблицы activities
CREATE TABLE activities (
                            id SERIAL PRIMARY KEY,                 -- Уникальный идентификатор
                            type VARCHAR(255) NOT NULL,            -- Тип активности
                            time TIMESTAMP NOT NULL,               -- Время активности
                            created_at TIMESTAMP DEFAULT NOW(),    -- Время создания записи
                            duration FLOAT,                        -- Длительность активности (в минутах), может быть NULL
                            user_id INT NOT NULL                   -- ID пользователя
);

-- Создание функции update_duration
CREATE OR REPLACE FUNCTION update_duration() RETURNS TRIGGER
    LANGUAGE plpgsql
AS $$
DECLARE
    previous_activity RECORD;
    next_activity RECORD;
BEGIN
    -- Находим предыдущую активность, которая предшествует новой
    SELECT * INTO previous_activity
    FROM activities
    WHERE user_id = NEW.user_id AND time < NEW.time
    ORDER BY time DESC
    LIMIT 1;

    -- Если найдена предыдущая активность, обновляем её длительность
    IF FOUND THEN
        RAISE NOTICE 'Updating duration for previous activity ID: %, Previous time: %, New time: %',
            previous_activity.id,
            previous_activity.time,
            NEW.time;

        UPDATE activities
        SET duration = EXTRACT(EPOCH FROM (NEW.time - previous_activity.time)) / 60
        WHERE id = previous_activity.id;

        -- Проверяем, если у новой активности есть следующая
        SELECT * INTO next_activity
        FROM activities
        WHERE user_id = NEW.user_id AND time > NEW.time
        ORDER BY time ASC
        LIMIT 1;

        -- Если следующая активность найдена, обновляем её длительность
        IF FOUND THEN
            RAISE NOTICE 'Updating duration for next activity ID: %, Next time: %, Previous time: %',
                next_activity.id,
                next_activity.time,
                NEW.time;

            UPDATE activities
            SET duration = EXTRACT(EPOCH FROM (next_activity.time - NEW.time)) / 60
            WHERE id = NEW.id;
        END IF;
    ELSE
        RAISE NOTICE 'No previous activity found for user ID: % at time: %',
            NEW.user_id,
            NEW.time;
    END IF;

    RETURN NEW;
END;
$$;

-- Назначение владельца функции (если требуется)
    ALTER FUNCTION update_duration() OWNER TO postgres;

-- Создание триггера set_duration_on_insert
CREATE TRIGGER set_duration_on_insert
    AFTER INSERT ON activities
    FOR EACH ROW
EXECUTE FUNCTION update_duration();
