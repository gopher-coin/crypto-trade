CREATE TABLE IF NOT EXISTS ohlc (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(20) NOT NULL,
    open NUMERIC(18, 8) NOT NULL,
    high NUMERIC(18, 8) NOT NULL,
    low NUMERIC(18, 8) NOT NULL,
    close NUMERIC(18, 8) NOT NULL,
    volume NUMERIC(18, 8),
    timestamp BIGINT NOT NULL,
    UNIQUE (symbol, timestamp)
);

CREATE TABLE IF NOT EXISTS trades (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(20) NOT NULL,
    side VARCHAR(10) NOT NULL,
    quantity NUMERIC(18, 8) NOT NULL,
    price NUMERIC(18, 8) NOT NULL,
    timestamp BIGINT NOT NULL,
    total NUMERIC(18, 8),
    UNIQUE (symbol, timestamp, side)
);