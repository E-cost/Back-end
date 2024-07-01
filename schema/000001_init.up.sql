START TRANSACTION;

CREATE EXTENSION IF NOT EXISTS pgcrypto;
COMMIT;

CREATE TABLE IF NOT EXISTS contacts
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(40) NOT NULL,
    surname VARCHAR(40),
    email VARCHAR(254) NOT NULL,
    contact_phone VARCHAR(15) NOT NULL,
    message VARCHAR(180),
    ip_address INET,
    port VARCHAR(5),
    proxy_chain TEXT[],
    created_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'Europe/Moscow'),
    updated_at TIMESTAMP
);
COMMIT;

CREATE TABLE IF NOT EXISTS confirmations
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    contact_id UUID NOT NULL,
    secret_code VARCHAR(50) NOT NULL,
    is_used BOOLEAN NOT NULL DEFAULT false,
    ip_address CIDR,
    port VARCHAR(5),
    proxy_chain TEXT[],
    created_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'Europe/Moscow'),
    confirmed_at TIMESTAMP,
    expired_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'Europe/Moscow' + INTERVAL '8 hours'),
    CONSTRAINT fk_confirmations_contacts
        FOREIGN KEY (contact_id)
            REFERENCES contacts(id)
            ON DELETE CASCADE
);
COMMIT;

CREATE TABLE IF NOT EXISTS orders
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    short_id VARCHAR(12) DEFAULT 'ODR-' || SUBSTRING(md5(gen_random_uuid()::text), 1, 8),
    contact_id UUID NOT NULL,
    items JSONB NOT NULL,
    total_amount_byn NUMERIC(10,2),
    total_amount_usd NUMERIC(10,2),
    is_paid BOOLEAN NOT NULL DEFAULT false,
    ip_address CIDR,
    port VARCHAR(5),
    proxy_chain TEXT[],
    created_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'Europe/Moscow'),
    expired_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'Europe/Moscow' + INTERVAL '14 days'),
    CONSTRAINT fk_orders_contacts
        FOREIGN KEY (contact_id)
            REFERENCES contacts(id)
            ON DELETE CASCADE
);
COMMIT;

CREATE TABLE IF NOT EXISTS seafood
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    article VARCHAR(20) DEFAULT 'ECBY-' || SUBSTRING(md5(gen_random_uuid()::text), 1, 8),
    category VARCHAR(50) NOT NULL,
    product VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    country_id VARCHAR(5) NOT NULL,
    net_weight VARCHAR(50) NOT NULL,
    composition VARCHAR(500),
    food_value VARCHAR(400),
    supplements VARCHAR(400),
    vitamins VARCHAR(300),
    energy_value VARCHAR(200),
    description VARCHAR(500),
    recommendation VARCHAR(500),
    shelf_life VARCHAR(500),
    expiration_date VARCHAR(100),
    price_byn NUMERIC(10,2),
    price_usd NUMERIC(10,2),
    created_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'Europe/Moscow'),
    updated_at TIMESTAMP
);
COMMIT;

INSERT INTO seafood
(category, product, name, country_id, net_weight, composition, food_value, vitamins, energy_value, description, recommendation, shelf_life, expiration_date, price_byn, price_usd)
VALUES
(
    'seafood', 'caviar', 'Горбуша с/м, б/к', 'RUS', '125г',
    'Икра (зерно) дальневосточных лососевых рыб, соль поваренная, глицерин пищевой, масло подсолнечное рафинированное.',
    '100г продукта (средняя): белок - 32г, жир - 15г',
    'В1 - 0.2 мг, В2 - 0.11 мг, РР - 1.2 мг', '263ккал/1099кДж',
    'Горбуша - наиболее распространенный представитель семейства лососевых, именно поэтому икра этой рыбы является самой популярной и привычной для большинства. Обладает по истине классическим вкусом, c легким оттенком горечи.',
    'Икру перед употреблением рекомендуется размораживать при температуре от 4°C до 15°С',
    'После вскрытия упаковки возможно кратковременное (не более 36 часов) хранение продукта при температуре от 0 до 6°С в условиях домашнего холодильника. Хранить при температуре не выше минус 18°С. Упаковано под вакуумом.',
    '12 месяцев с даты изготовления.',
    42.99,
    13.5
);

INSERT INTO seafood
(category, product, name, country_id, net_weight, composition, food_value, vitamins, energy_value, description, recommendation, shelf_life, expiration_date, price_byn, price_usd)
VALUES
(
    'seafood', 'caviar', 'Горбуша с/м, б/к', 'RUS', '240г',
    'Икра (зерно) дальневосточных лососевых рыб, соль поваренная, глицерин пищевой, масло подсолнечное рафинированное.',
    '100г продукта (средняя): белок - 32г, жир - 15г',
    'В1 - 0.2 мг, В2 - 0.11 мг, РР - 1.2 мг', '263ккал/1099кДж',
    'Горбуша - наиболее распространенный представитель семейства лососевых, именно поэтому икра этой рыбы является самой популярной и привычной для большинства. Обладает по истине классическим вкусом, c легким оттенком горечи.',
    'Икру перед употреблением рекомендуется размораживать при температуре от 4°C до 15°С',
    'После вскрытия упаковки возможно кратковременное (не более 36 часов) хранение продукта при температуре от 0 до 6°С в условиях домашнего холодильника. Хранить при температуре не выше минус 18°С. Упаковано под вакуумом.',
    '12 месяцев с даты изготовления.',
    42.99,
    13.5
);

INSERT INTO seafood
(category, product, name, country_id, net_weight, composition, food_value, supplements, vitamins, energy_value, description, recommendation, shelf_life, expiration_date, price_byn, price_usd)
VALUES
(
    'seafood', 'caviar', 'Горбуша соленая', 'RUS', '125г',
    'Икра (зерно) дальневосточных лососевых рыб, соль поваренная, глицерин пищевой, масло подсолнечное рафинированное.',
    '100г продукта (средняя): белок - 32г, жир - 15г',
    'Глицерин пищевой Е422; консерванты: E211, Е200.',
    'В1 - 0.2 мг, В2 - 0.11 мг, РР - 1.2 мг', '263ккал/1099кДж',
    'Горбуша - наиболее распространенный представитель семейства лососевых, именно поэтому икра этой рыбы является самой популярной и привычной для большинства. Обладает по истине классическим вкусом, c легким оттенком горечи.',
    'Икру перед употреблением рекомендуется размораживать при температуре от 4°C до 15°С',
    'После вскрытия упаковки возможно кратковременное (не более 36 часов) хранение продукта при температуре от 0 до 6°С в условиях домашнего холодильника. Хранить при температуре не выше минус 18°С. Упаковано под вакуумом.',
    '12 месяцев с даты изготовления.',
    42.99,
    13.5
);

INSERT INTO seafood
(category, product, name, country_id, net_weight, composition, food_value, supplements, vitamins, energy_value, description, recommendation, shelf_life, expiration_date, price_byn, price_usd)
VALUES
(
    'seafood', 'caviar', 'Горбуша соленая', 'RUS', '240г',
    'Икра (зерно) дальневосточных лососевых рыб, соль поваренная, глицерин пищевой, масло подсолнечное рафинированное.',
    '100г продукта (средняя): белок - 32г, жир - 15г',
    'Глицерин пищевой Е422; консерванты: E211, Е200.',
    'В1 - 0.2 мг, В2 - 0.11 мг, РР - 1.2 мг', '263ккал/1099кДж',
    'Горбуша - наиболее распространенный представитель семейства лососевых, именно поэтому икра этой рыбы является самой популярной и привычной для большинства. Обладает по истине классическим вкусом, c легким оттенком горечи.',
    'Икру перед употреблением рекомендуется размораживать при температуре от 4°C до 15°С',
    'После вскрытия упаковки возможно кратковременное (не более 36 часов) хранение продукта при температуре от 0 до 6°С в условиях домашнего холодильника. Хранить при температуре не выше минус 18°С. Упаковано под вакуумом.',
    '12 месяцев с даты изготовления.',
    42.99,
    13.5
);

INSERT INTO seafood
(category, product, name, country_id, net_weight, composition, food_value, vitamins, energy_value, description, recommendation, shelf_life, expiration_date, price_byn, price_usd)
VALUES
(
    'seafood', 'caviar', 'Кета с/м, б/к', 'RUS', '125г',
    'Икра (зерно) дальневосточных лососевых рыб, соль поваренная, глицерин пищевой, масло подсолнечное рафинированное.',
    '100г продукта (средняя): белок - 32г, жир - 15г',
    'В1 - 0.2 мг, В2 - 0.11 мг, РР - 1.2 мг', '263ккал/1099кДж',
    'Настоящий деликатес, который высоко ценится среди других видов, по размеру уступает только икре чавычи, занесенной в Красную книгу. Не даром икру кеты называют царской, за ее янтарный блеск, крупные зерна и нежный сливочный вкус.',
    'Икру перед употреблением рекомендуется размораживать при температуре от 4°C до 15°С',
    'После вскрытия упаковки возможно кратковременное (не более 36 часов) хранение продукта при температуре от 0 до 6°С в условиях домашнего холодильника. Хранить при температуре не выше минус 18°С. Упаковано под вакуумом.',
    '12 месяцев с даты изготовления.',
    42.99,
    13.5
);

INSERT INTO seafood
(category, product, name, country_id, net_weight, composition, food_value, vitamins, energy_value, description, recommendation, shelf_life, expiration_date, price_byn, price_usd)
VALUES
(
    'seafood', 'caviar', 'Кета с/м, б/к', 'RUS', '240г',
    'Икра (зерно) дальневосточных лососевых рыб, соль поваренная, глицерин пищевой, масло подсолнечное рафинированное.',
    '100г продукта (средняя): белок - 32г, жир - 15г',
    'В1 - 0.2 мг, В2 - 0.11 мг, РР - 1.2 мг', '263ккал/1099кДж',
    'Настоящий деликатес, который высоко ценится среди других видов, по размеру уступает только икре чавычи, занесенной в Красную книгу. Не даром икру кеты называют царской, за ее янтарный блеск, крупные зерна и нежный сливочный вкус.',
    'Икру перед употреблением рекомендуется размораживать при температуре от 4°C до 15°С',
    'После вскрытия упаковки возможно кратковременное (не более 36 часов) хранение продукта при температуре от 0 до 6°С в условиях домашнего холодильника. Хранить при температуре не выше минус 18°С. Упаковано под вакуумом.',
    '12 месяцев с даты изготовления.',
    42.99,
    13.5
);

INSERT INTO seafood
(category, product, name, country_id, net_weight, composition, food_value, supplements, vitamins, energy_value, description, recommendation, shelf_life, expiration_date, price_byn, price_usd)
VALUES
(
    'seafood', 'caviar', 'Кета соленая', 'RUS', '125г',
    'Икра (зерно) дальневосточных лососевых рыб, соль поваренная, глицерин пищевой, масло подсолнечное рафинированное.',
    '100г продукта (средняя): белок - 32г, жир - 15г',
    'Глицерин пищевой Е422; консерванты: E211, Е200.',
    'В1 - 0.2 мг, В2 - 0.11 мг, РР - 1.2 мг', '263ккал/1099кДж',
    'Настоящий деликатес, который высоко ценится среди других видов, по размеру уступает только икре чавычи, занесенной в Красную книгу. Не даром икру кеты называют царской, за ее янтарный блеск, крупные зерна и нежный сливочный вкус.',
    'Икру перед употреблением рекомендуется размораживать при температуре от 4°C до 15°С',
    'После вскрытия упаковки возможно кратковременное (не более 36 часов) хранение продукта при температуре от 0 до 6°С в условиях домашнего холодильника. Хранить при температуре не выше минус 18°С. Упаковано под вакуумом.',
    '12 месяцев с даты изготовления.',
    42.99,
    13.5
);

INSERT INTO seafood
(category, product, name, country_id, net_weight, composition, food_value, supplements, vitamins, energy_value, description, recommendation, shelf_life, expiration_date, price_byn, price_usd)
VALUES
(
    'seafood', 'caviar', 'Кета соленая', 'RUS', '240г',
    'Икра (зерно) дальневосточных лососевых рыб, соль поваренная, глицерин пищевой, масло подсолнечное рафинированное.',
    '100г продукта (средняя): белок - 32г, жир - 15г',
    'Глицерин пищевой Е422; консерванты: E211, Е200.',
    'В1 - 0.2 мг, В2 - 0.11 мг, РР - 1.2 мг', '263ккал/1099кДж',
    'Настоящий деликатес, который высоко ценится среди других видов, по размеру уступает только икре чавычи, занесенной в Красную книгу. Не даром икру кеты называют царской, за ее янтарный блеск, крупные зерна и нежный сливочный вкус.',
    'Икру перед употреблением рекомендуется размораживать при температуре от 4°C до 15°С',
    'После вскрытия упаковки возможно кратковременное (не более 36 часов) хранение продукта при температуре от 0 до 6°С в условиях домашнего холодильника. Хранить при температуре не выше минус 18°С. Упаковано под вакуумом.',
    '12 месяцев с даты изготовления.',
    42.99,
    13.5
);


INSERT INTO seafood
(category, product, name, country_id, net_weight, composition, food_value, vitamins, energy_value, description, recommendation, shelf_life, expiration_date, price_byn, price_usd)
VALUES
(
    'seafood', 'caviar', 'Нерка с/м, б/к', 'RUS', '125г',
    'Икра (зерно) дальневосточных лососевых рыб, соль поваренная, глицерин пищевой, масло подсолнечное рафинированное.',
    '100г продукта (средняя): белок - 32г, жир - 15г',
    'В1 - 0.2 мг, В2 - 0.11 мг, РР - 1.2 мг', '263ккал/1099кДж',
    'Икра нерки является самой дефицитной из всех видов, обладает насыщенными вкусовыми красками, с выраженным оттенком горечи, тем не менее она так нежна, что в буквальном смысле слова, тает во рту.',
    'Икру перед употреблением рекомендуется размораживать при температуре от 4°C до 15°С',
    'После вскрытия упаковки возможно кратковременное (не более 36 часов) хранение продукта при температуре от 0 до 6°С в условиях домашнего холодильника. Хранить при температуре не выше минус 18°С. Упаковано под вакуумом.',
    '12 месяцев с даты изготовления.',
    42.99,
    13.5
);

INSERT INTO seafood
(category, product, name, country_id, net_weight, composition, food_value, supplements, vitamins, energy_value, description, recommendation, shelf_life, expiration_date, price_byn, price_usd)
VALUES
(
    'seafood', 'caviar', 'Нерка соленая', 'RUS', '125г',
    'Икра (зерно) дальневосточных лососевых рыб, соль поваренная, глицерин пищевой, масло подсолнечное рафинированное.',
    '100г продукта (средняя): белок - 32г, жир - 15г',
    'Глицерин пищевой Е422; консерванты: E211, Е200.',
    'В1 - 0.2 мг, В2 - 0.11 мг, РР - 1.2 мг', '263ккал/1099кДж',
    'Икра нерки является самой дефицитной из всех видов, обладает насыщенными вкусовыми красками, с выраженным оттенком горечи, тем не менее она так нежна, что в буквальном смысле слова, тает во рту.',
    'Икру перед употреблением рекомендуется размораживать при температуре от 4°C до 15°С',
    'После вскрытия упаковки возможно кратковременное (не более 36 часов) хранение продукта при температуре от 0 до 6°С в условиях домашнего холодильника. Хранить при температуре не выше минус 18°С. Упаковано под вакуумом.',
    '12 месяцев с даты изготовления.',
    42.99,
    13.5
);