package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/luc-wallace/idolvault/internal/models"
)

func (c *Conn) GetGroups(ctx context.Context) ([]*models.Group, error) {
	rows, err := c.conn.Query(ctx, "SELECT * FROM groups ORDER BY popularity DESC")
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.Group])
}

func (c *Conn) GetGroup(ctx context.Context, groupName string) (*models.Group, error) {
	row, err := c.conn.Query(ctx, "SELECT * FROM groups WHERE LOWER(name) = LOWER($1)", groupName)
	if err != nil {
		return nil, err
	}
	group, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByName[models.Group])
	if errors.Is(err, pgx.ErrNoRows) {
		err = nil
	}
	return group, err
}

func (c *Conn) UpdateGroupStats(ctx context.Context, group *models.Group) error {
	_, err := c.conn.Exec(ctx, `
		UPDATE groups
		SET popularity = $1, followers = $2, image_url = $3, genres = $4, updated_at = $5
		WHERE name = $6
	`, group.Popularity, group.Followers, group.ImageURL, group.Genres, group.UpdatedAt, group.Name)
	return err
}
