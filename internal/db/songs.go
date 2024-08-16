package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/luc-wallace/idolvault/internal/models"
)

func (c *Conn) CreateSong(ctx context.Context, song *models.Song) error {
	_, err := c.conn.Exec(ctx, "INSERT INTO songs VALUES ($1, $2, $3, $4, $5)", song.SpotifyID, song.Name, song.GroupName, song.Popularity, song.AlbumImageURL)
	return err
}

func (c *Conn) DeleteGroupSongs(ctx context.Context, groupName string) error {
	_, err := c.conn.Exec(ctx, "DELETE FROM songs WHERE group_name = $1", groupName)
	return err
}

func (c *Conn) GetGroupSongs(ctx context.Context, groupName string) ([]*models.Song, error) {
	rows, err := c.conn.Query(ctx, "SELECT * FROM songs WHERE LOWER(group_name) = LOWER($1) ORDER BY popularity DESC", groupName)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.Song])
}
