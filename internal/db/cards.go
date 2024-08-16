package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/luc-wallace/idolvault/internal/models"
)

func (c *Conn) GetCollectionCards(ctx context.Context, groupName string, collectionName string) ([]*models.CardWithOwnershipState, error) {
	rows, err := c.conn.Query(
		ctx,
		"SELECT * FROM cards WHERE LOWER(group_name) = LOWER($1) AND LOWER(collection_name) = LOWER($2)",
		groupName,
		collectionName,
	)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[models.CardWithOwnershipState])
}

func (c *Conn) GetCardOwnershipState(ctx context.Context, username string, cardID int) (bool, error) {
	var exists bool
	if err := c.conn.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM user_cards WHERE username = $1 AND card_id = $2)", username, cardID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (c *Conn) SetCardOwned(ctx context.Context, username string, cardID int) error {
	_, err := c.conn.Exec(ctx, "INSERT INTO user_cards VALUES ($1, $2)", username, cardID)
	return err
}

func (c *Conn) SetCardUnowned(ctx context.Context, username string, cardID int) error {
	_, err := c.conn.Exec(ctx, "DELETE FROM user_cards WHERE username = $1 AND card_id = $2", username, cardID)
	return err
}

// idolName is optional
func (c *Conn) GetCollectionCardsWithOwnershipState(ctx context.Context, username string, groupName string, collectionName string, idolName string) ([]*models.CardWithOwnershipState, error) {
	sql := `SELECT
			cards.*,
			(SELECT EXISTS (SELECT 1 FROM user_cards WHERE username = $1 AND card_id = cards.id)) AS owned
		FROM cards
		WHERE
			LOWER(group_name) = LOWER($2)
			AND
			LOWER(collection_name) = LOWER($3)`

	args := []any{username, groupName, collectionName}
	var rows pgx.Rows
	var err error

	if idolName != "" {
		sql += ` AND LOWER(idol_name) = LOWER($4)`
		args = append(args, idolName)
	}

	rows, err = c.conn.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.CardWithOwnershipState])
}

func (c *Conn) GetCard(ctx context.Context, id int) (*models.Card, error) {
	row, err := c.conn.Query(ctx, "SELECT * FROM cards WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[models.Card])
}
