CREATE TABLE pvz
(
    id                UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    registration_date TIMESTAMPTZ NOT NULL DEFAULT now(),
    city              VARCHAR(255) CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань'))
);

CREATE TABLE receptions
(
    id        UUID PRIMARY KEY                                                 DEFAULT gen_random_uuid(),
    date_time TIMESTAMPTZ                                             NOT NULL DEFAULT now(),
    pvz_id    UUID                                                    NOT NULL,
    status    VARCHAR(255) CHECK (status IN ('in_progress', 'close')) NOT NULL,
    FOREIGN KEY (pvz_id) REFERENCES pvz (id)
);

CREATE TABLE products
(
    id           UUID PRIMARY KEY,
    date_time    TIMESTAMPTZ DEFAULT now()                                      NOT NULL,
    type         VARCHAR(50) CHECK (type IN ('электроника', 'одежда', 'обувь')) NOT NULL,
    reception_id UUID REFERENCES receptions (id)
);
CREATE INDEX idx_pvz_id ON receptions (pvz_id);

CREATE TABLE users
(
    id       UUID PRIMARY KEY,
    email    VARCHAR(255) UNIQUE                                   NOT NULL,
    password VARCHAR(255)                                          NOT NULL,
    role     VARCHAR(50) CHECK (role IN ('employee', 'moderator')) NOT NULL
);
CREATE INDEX idx_role ON users (role);