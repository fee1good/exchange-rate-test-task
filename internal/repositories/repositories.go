package repositories

import (
	"github.com/fee1good/exchange-rate-test-task/internal/repositories/rates"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Container struct {
	Rates *rates.Repository
}

func New(pool *pgxpool.Pool) *Container {
	return &Container{
		Rates: rates.NewRepository(pool),
	}
}
