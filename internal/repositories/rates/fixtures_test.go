package rates

type cleanupFixture struct{}

func (f cleanupFixture) GetSql() []string {
	return []string{
		`TRUNCATE TABLE rates RESTART IDENTITY CASCADE;`,
	}
}

type insertRatePairsFixture struct{}

func (f insertRatePairsFixture) GetSql() []string {
	return []string{
		`INSERT INTO rates (crypto_symbol, fiat_symbol, raw_data, display_data) VALUES ('ETH', 'USD', '{"CHANGE24HOUR":1}', '{"CHANGE24HOUR":"1"}');`,
		`INSERT INTO rates (crypto_symbol, fiat_symbol, raw_data, display_data) VALUES ('ETH', 'EUR', '{"CHANGE24HOUR":2}', '{"CHANGE24HOUR":"2"}');`,
	}
}
