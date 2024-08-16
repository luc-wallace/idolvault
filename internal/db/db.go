package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Conn struct {
	conn *pgxpool.Pool
}

func Connect(ctx context.Context, uri string) (*Conn, error) {
	conn, err := pgxpool.New(ctx, uri)
	if err != nil {
		return nil, err
	}
	return &Conn{conn: conn}, nil
}

func (c *Conn) Ping(ctx context.Context) error {
	return c.conn.Ping(ctx)
}

func (c *Conn) Pool() *pgxpool.Pool {
	return c.conn
}

func (c *Conn) Close(ctx context.Context) {
	c.conn.Close()
}
