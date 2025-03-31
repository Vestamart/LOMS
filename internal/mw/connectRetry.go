package mw

import (
	"context"
	"github.com/jackc/pgx/v5"
	"time"
)

func ConnectWithRetry(ctx context.Context, dsn string, maxAttempts int, delay time.Duration) (*pgx.Conn, error) {
	var conn *pgx.Conn
	var err error
	for i := 0; i < maxAttempts; i++ {
		conn, err = pgx.Connect(ctx, dsn)
		if err == nil {
			return conn, nil
		}
		time.Sleep(delay)
	}
	return nil, err
}
