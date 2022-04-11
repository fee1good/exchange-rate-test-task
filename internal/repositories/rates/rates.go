package rates

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/multierr"

	"github.com/fee1good/exchange-rate-test-task/internal/entities"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) UploadPairs(ctx context.Context, pairs []entities.Pair) error {
	const batchSize = 100

	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	var startIndex, endIndex = 0, batchSize
	for startIndex < len(pairs) {
		if endIndex > len(pairs) {
			endIndex = len(pairs)
		}

		batch := pairs[startIndex:endIndex]
		if err := r.uploadPairs(ctx, tx, batch); err != nil {
			return multierr.Append(err, tx.Rollback(ctx))
		}

		startIndex = endIndex
		endIndex += batchSize
	}

	return tx.Commit(ctx)
}

func (r *Repository) uploadPairs(ctx context.Context, tx pgx.Tx, pairs []entities.Pair) error {
	batch := pgx.Batch{}
	for _, pair := range pairs {
		batch.Queue(
			upsertPairRatesQuery,
			pair.CryptoSymbol,
			pair.FiatSymbol,
			pair.RawRateValues,
			pair.DisplayRateValues,
		)
	}

	batchResult := tx.SendBatch(ctx, &batch)
	for i := 0; i < batch.Len(); i++ {
		if _, err := batchResult.Exec(); err != nil {
			return err
		}
	}

	return batchResult.Close()
}

func (r *Repository) GetPairs(ctx context.Context, cryptoSymbols, fiatSymbols []string) ([]entities.Pair, error) {
	conn, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, getPairRatesQuery, cryptoSymbols, fiatSymbols)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pairs []entities.Pair
	for rows.Next() {
		var pair entities.Pair
		err := rows.Scan(
			&pair.CryptoSymbol,
			&pair.FiatSymbol,
			&pair.RawRateValues,
			&pair.DisplayRateValues,
		)
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, pair)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pairs, nil
}
