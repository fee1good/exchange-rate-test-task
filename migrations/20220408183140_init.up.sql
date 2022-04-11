CREATE TABLE IF NOT EXISTS rates (
    crypto_symbol text not null,
    fiat_symbol text not null,
    raw_data json not null,
    display_data json not null,
    updated_at timestamp default now(),
    CONSTRAINT rates_primary_key PRIMARY KEY (
        crypto_symbol, fiat_symbol
    )
);