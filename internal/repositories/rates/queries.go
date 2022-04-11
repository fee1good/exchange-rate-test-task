package rates

const upsertPairRatesQuery = `
INSERT INTO rates 
	(crypto_symbol, fiat_symbol, raw_data, display_data)
VALUES 
	($1, $2, $3, $4)
ON CONFLICT ON CONSTRAINT rates_primary_key
DO UPDATE SET
	raw_data = EXCLUDED.raw_data,
	display_data = EXCLUDED.display_data,
	updated_at = NOW()
`

const getPairRatesQuery = `
SELECT
	crypto_symbol,
	fiat_symbol,
	raw_data,
	display_data
FROM rates
WHERE crypto_symbol IN ($1) AND fiat_symbol IN ($2)
`
