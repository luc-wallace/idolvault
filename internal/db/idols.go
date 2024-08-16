package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/luc-wallace/idolvault/internal/models"
)

func (c *Conn) GetGroupIdols(ctx context.Context, groupName string) ([]*models.Idol, error) {
	rows, err := c.conn.Query(ctx, "SELECT * FROM idols WHERE LOWER(group_name) = LOWER($1)", groupName)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.Idol])
}

func (c *Conn) GetGroupIdol(ctx context.Context, groupName string, idolName string) (*models.Idol, error) {
	row, err := c.conn.Query(ctx, "SELECT * FROM idols WHERE LOWER(group_name) = LOWER($1) AND LOWER(stage_name) = LOWER($2)", groupName, idolName)
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[models.Idol])
}

func (c *Conn) GetGroupIdolWithBias(ctx context.Context, groupName string, idolName string, username string) (*models.IdolWithBias, error) {
	row, err := c.conn.Query(ctx, `
	SELECT idols.*,
	(SELECT EXISTS(
		SELECT 1 FROM biases WHERE username = $1
		AND biases.idol_name = idols.stage_name
		AND biases.group_name = idols.group_name
	)) AS bias
	FROM idols WHERE LOWER(group_name) = LOWER($2) AND LOWER(stage_name) = LOWER($3)`, username, groupName, idolName)
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[models.IdolWithBias])
}
