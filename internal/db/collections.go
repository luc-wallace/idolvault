package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/luc-wallace/idolvault/internal/models"
)

func (c *Conn) GetGroupCollections(ctx context.Context, groupName string, asc bool) ([]*models.Collection, error) {
	query := "SELECT * FROM collections WHERE LOWER(group_name) = LOWER($1) ORDER BY release_date "
	if asc {
		query += "ASC"
	} else {
		query += "DESC"
	}

	rows, err := c.conn.Query(ctx, query, groupName)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.Collection])
}

func (c *Conn) GetGroupCollectionsWithCardCounts(ctx context.Context, groupName string, asc bool) ([]*models.CollectionWithCardCount, error) {
	query := `SELECT collections.*, 
			(SELECT COUNT(*) FROM cards WHERE
				cards.collection_name = collections.name AND
				cards.group_name = collections.group_name
			) AS card_count FROM collections WHERE LOWER(group_name) = LOWER($1) ORDER BY release_date`
	if asc {
		query += " ASC"
	} else {
		query += " DESC"
	}

	rows, err := c.conn.Query(ctx, query, groupName)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.CollectionWithCardCount])
}
