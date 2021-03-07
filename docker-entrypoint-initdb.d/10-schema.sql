CREATE TABLE IF NOT EXISTS clients (
    id BIGSERIAL PRIMARY KEY,
    login VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(512) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    passport VARCHAR(10) NOT NULL CHECK (length(passport) = 10),
    birthday DATE NOT NULL,
    status VARCHAR(8) NOT NULL DEFAULT 'INACTIVE' CHECK (status IN ('INACTIVE', 'ACTIVE')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cards (
    id BIGSERIAL PRIMARY KEY,
    number VARCHAR(20) NOT NULL UNIQUE,
    balance BIGINT NOT NULL DEFAULT 0,
    issuer VARCHAR(10) NOT NULL CHECK (issuer IN ('Visa', 'MasterCard', 'MIR')),
    holder VARCHAR(80) NOT NULL, -- имя держателя на карте
    owner_id BIGINT NOT NULL REFERENCES clients(id),
    status VARCHAR(8) NOT NULL DEFAULT 'INACTIVE' CHECK (status IN ('INACTIVE', 'ACTIVE')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS icons (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL UNIQUE,
    uri VARCHAR(500) NOT NULL
);

CREATE TABLE IF NOT EXISTS mccs (
    id VARCHAR(4) PRIMARY KEY,
    text VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    from_id BIGINT REFERENCES cards(id),
    to_id BIGINT REFERENCES cards(id),
    sum BIGINT NOT NULL,
    mcc_id VARCHAR(4) REFERENCES mccs(id),
    icon_id BIGINT REFERENCES icons(id),
    description VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);