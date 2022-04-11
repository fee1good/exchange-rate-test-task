package rates

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"

	"github.com/fee1good/exchange-rate-test-task/internal/config"
	"github.com/fee1good/exchange-rate-test-task/internal/entities"
	"github.com/fee1good/exchange-rate-test-task/internal/system/connections/pg"
	"github.com/fee1good/exchange-rate-test-task/test/fixtures"
)

func getTestRepoAndPool(t *testing.T, ctx context.Context) (*Repository, *pgxpool.Pool) {
	cfg, err := config.New()
	assert.NoError(t, err)

	pool, err := pg.NewPool(ctx, cfg.PGDSN)
	assert.NoError(t, err)

	return NewRepository(pool), pool
}

func sqlBoolCheck(t *testing.T, ctx context.Context, pool *pgxpool.Pool, query string) bool {
	conn, err := pool.Acquire(ctx)
	assert.NoError(t, err)
	defer conn.Release()

	row := conn.QueryRow(context.Background(), query)
	var answer bool
	assert.NoError(t, row.Scan(&answer))
	return answer
}

func execQuery(t *testing.T, ctx context.Context, pool *pgxpool.Pool, query string) {
	conn, err := pool.Acquire(ctx)
	assert.NoError(t, err)
	defer conn.Release()

	_, err = conn.Exec(context.Background(), query)
	assert.NoError(t, err)
}

func TestRepository_UploadPairs(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	bigCountPairs := make([]entities.Pair, 0, 150)
	for i := 0; i < 150; i++ {
		bigCountPairs = append(bigCountPairs, entities.Pair{
			CryptoSymbol:      fmt.Sprintf("cry_%d", i),
			FiatSymbol:        fmt.Sprintf("fia_%d", i),
			RawRateValues:     entities.RateValues[float64]{},
			DisplayRateValues: entities.RateValues[string]{},
		})
	}

	cases := []struct {
		name       string
		insertData []entities.Pair
		sqlCheck   string
		isError    bool
		errMessage string
	}{
		{
			name: "OK insert pairs",
			insertData: []entities.Pair{
				{"ETH", "USD", entities.RateValues[float64]{}, entities.RateValues[string]{}},
				{"ETH", "EUR", entities.RateValues[float64]{}, entities.RateValues[string]{}},
			},
			sqlCheck: `SELECT count = 2 FROM (SELECT count(*) FROM rates) AS p;`,
		},
		{
			name: "OK insert and update pair",
			insertData: []entities.Pair{
				{"ETH", "USD", entities.RateValues[float64]{}, entities.RateValues[string]{}},
				{"ETH", "USD", entities.RateValues[float64]{Change24Hour: 1}, entities.RateValues[string]{Change24Hour: "1"}},
			},
			sqlCheck: `SELECT count = 1 FROM
(SELECT count(*) FROM rates WHERE crypto_symbol = 'ETH' AND fiat_symbol = 'USD'
AND (raw_data->>'CHANGE24HOUR')::float = 1 AND (display_data->>'CHANGE24HOUR')::text = '1') AS p;`,
		},
		{
			name:       "OK insert a lot of pairs",
			insertData: bigCountPairs,
			sqlCheck:   `SELECT count = 150 FROM (SELECT count(*) FROM rates) AS p;`,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			var (
				ctx              = context.Background()
				repository, pool = getTestRepoAndPool(t, ctx)
			)
			defer pool.Close()

			fixtures.ExecuteFixture(t, ctx, pool, cleanupFixture{})

			err := repository.UploadPairs(ctx, testCase.insertData)
			if testCase.isError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testCase.errMessage)
			} else {
				assert.NoError(t, err)
			}

			assert.True(t, sqlBoolCheck(t, ctx, pool, testCase.sqlCheck))
		})
	}
}

//func TestRepository_GetPairs(t *testing.T) {
//	if testing.Short() {
//		t.Skip()
//	}
//
//	cases := []struct {
//		name                       string
//		cryptoSymbols, fiatSymbols []string
//		result                     []entities.Pair
//		isError                    bool
//		errMessage                 string
//	}{
//		{
//			name:          "OK get exists pairs",
//			cryptoSymbols: []string{"BTC", "ETH"},
//			fiatSymbols:   []string{"USD", "EUR"},
//			result: []entities.Pair{
//				{"ETH", "USD", entities.RateValues[float64]{Change24Hour: 1}, entities.RateValues[string]{Change24Hour: "1"}},
//				{"ETH", "EUR", entities.RateValues[float64]{Change24Hour: 2}, entities.RateValues[string]{Change24Hour: "2"}},
//			},
//		},
//	}
//
//	for _, testCase := range cases {
//		t.Run(testCase.name, func(t *testing.T) {
//			var (
//				ctx              = context.Background()
//				repository, pool = getTestRepoAndPool(t, ctx)
//			)
//			defer pool.Close()
//
//			fixtures.ExecuteFixture(t, ctx, pool, cleanupFixture{})
//			fixtures.ExecuteFixture(t, ctx, pool, insertRatePairsFixture{})
//
//			pairs, err := repository.GetPairs(ctx, testCase.cryptoSymbols, testCase.fiatSymbols)
//			if testCase.isError {
//				assert.Error(t, err)
//				assert.Contains(t, err.Error(), testCase.errMessage)
//			} else {
//				assert.NoError(t, err)
//			}
//
//			assert.Equal(t, testCase.result, pairs)
//		})
//	}
//}
