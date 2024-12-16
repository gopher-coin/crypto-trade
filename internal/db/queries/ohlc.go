package queries

import (
	"context"

	"github.com/gopher-coin/crypto-trade/pkg/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OHLCQueries struct {
	Pool *pgxpool.Pool
}

func (q *OHLCQueries) InsertOHLC(ctx context.Context, ohlc models.OHLC) error {
	query := `
		INSERT INTO ohlc (symbol, open, high, low, close, volume, timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (symbol, timestamp) DO NOTHING;
	`
	_, err := q.Pool.Exec(ctx, query,
		ohlc.Symbol,
		ohlc.Open,
		ohlc.High,
		ohlc.Low,
		ohlc.Close,
		ohlc.Volume,
		ohlc.Timestamp,
	)
	return err
}
