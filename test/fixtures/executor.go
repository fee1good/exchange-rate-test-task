package fixtures

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
)

type Fixture interface {
	GetSql() []string
}

func ExecuteFixture(t *testing.T, ctx context.Context, pool *pgxpool.Pool, fixture Fixture) {
	conn, err := pool.Acquire(ctx)
	assert.NoError(t, err)
	defer conn.Release()

	for _, query := range fixture.GetSql() {
		_, err := conn.Exec(ctx, query)
		assert.NoError(t, err)
	}
}
