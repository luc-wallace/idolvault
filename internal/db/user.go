package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/luc-wallace/idolvault/internal/models"
)

func (c *Conn) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	rows, err := c.conn.Query(ctx, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[models.User])
	if errors.Is(err, pgx.ErrNoRows) {
		err = nil
	}
	return user, err
}

func (c *Conn) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	rows, err := c.conn.Query(ctx, "SELECT * FROM users WHERE LOWER(username) = LOWER($1)", username)
	if err != nil {
		return nil, err
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[models.User])
	if errors.Is(err, pgx.ErrNoRows) {
		err = nil
	}
	return user, err
}

func (c *Conn) GetUserWithFollowing(ctx context.Context, username string, follower string) (*models.UserWithFollowers, error) {
	rows, err := c.conn.Query(ctx, `
	SELECT users.*,
	(SELECT EXISTS(SELECT 1 FROM followers WHERE LOWER(follower) = LOWER($1) AND LOWER(following) = LOWER($2))) AS is_following,
	(SELECT COUNT(*) FROM followers WHERE LOWER(following) = LOWER($3)) AS followers,
	(SELECT COUNT(*) FROM followers WHERE LOWER(follower) = LOWER($4)) AS following
	FROM users WHERE LOWER(username) = LOWER($5)`, follower, username, username, username, username)
	if err != nil {
		return nil, err
	}
	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[models.UserWithFollowers])
	if errors.Is(err, pgx.ErrNoRows) {
		err = nil
	}
	return user, err
}

func (c *Conn) CreateUser(ctx context.Context, user *models.User) error {
	_, err := c.conn.Exec(
		ctx,
		"INSERT INTO users VALUES ($1, $2, $3, $4, $5)",
		user.Email, user.Username, user.Provider, user.AvatarURL, user.Bio,
	)
	return err
}

func (c *Conn) GetUserCards(ctx context.Context, username string) ([]*models.Card, error) {
	rows, err := c.conn.Query(ctx, "SELECT cards.* FROM user_cards INNER JOIN cards ON user_cards.card_id = cards.id WHERE LOWER(username) = LOWER($1)", username)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.Card])
}

func (c *Conn) GetUserCardsWithFilter(ctx context.Context, username string, group string, collection string, idol string) ([]*models.Card, error) {
	sql := "SELECT cards.* FROM user_cards INNER JOIN cards ON user_cards.card_id = cards.id WHERE LOWER(username) = LOWER($1)"
	params := []any{username}
	if group != "" {
		params = append(params, group)
		sql += " AND group_name = $2"
	}
	if collection != "" {
		params = append(params, collection)
		sql += fmt.Sprintf(" AND collection_name = $%d", len(params))
	}
	if idol != "" {
		params = append(params, idol)
		sql += fmt.Sprintf(" AND idol_name = $%d", len(params))
	}

	fmt.Println(sql)
	fmt.Println(params)
	rows, err := c.conn.Query(ctx, sql, params...)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.Card])

}

func (c *Conn) GetUserFollowers(ctx context.Context, username string) ([]*models.User, error) {
	rows, err := c.conn.Query(ctx, `
		SELECT users.* FROM followers 
		INNER JOIN users ON LOWER(followers.follower) = LOWER(users.username)
		WHERE followers.following = $1
	`, username)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.User])
}

func (c *Conn) CreateUserFollower(ctx context.Context, follower string, following string) error {
	_, err := c.conn.Exec(ctx, "INSERT INTO followers (follower, following) VALUES ($1, $2)", follower, following)
	return err
}

func (c *Conn) DeleteUserFollower(ctx context.Context, follower string, following string) error {
	_, err := c.conn.Exec(ctx, "DELETE FROM followers WHERE follower = $1 AND following = $2", follower, following)
	return err
}

func (c *Conn) CreateUserBias(ctx context.Context, username string, idolName string, groupName string) error {
	_, err := c.conn.Exec(ctx, "INSERT INTO biases VALUES ($1, $2, $3)", username, idolName, groupName)
	return err
}

func (c *Conn) DeleteUserBias(ctx context.Context, username string, idolName string, groupName string) error {
	_, err := c.conn.Exec(ctx, "DELETE FROM biases WHERE LOWER(username) = LOWER($1) AND idol_name = $2 AND group_name = $3", username, idolName, groupName)
	return err
}

func (c *Conn) GetUserBiases(ctx context.Context, username string) ([]*models.Idol, error) {
	rows, err := c.conn.Query(ctx, `
	SELECT idols.* FROM biases
	INNER JOIN idols
		ON biases.group_name = idols.group_name
		AND biases.idol_name = idols.stage_name
	WHERE LOWER(username) = LOWER($1)`, username)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.Idol])
}
