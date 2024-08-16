package models

import (
	"encoding/gob"
	"time"
)

func init() {
	gob.Register(NewUser{})
	gob.Register(User{})
}

type Group struct {
	Name       string    `db:"name"`
	Country    string    `db:"country"`
	Fandom     string    `db:"fandom"`
	SpotifyID  string    `db:"spotify_id"`
	Popularity int       `db:"popularity"`
	Followers  int64     `db:"followers"`
	ImageURL   string    `db:"image_url"`
	Genres     []string  `db:"genres"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type Song struct {
	SpotifyID     string `db:"spotify_id"`
	Name          string `db:"name"`
	GroupName     string `db:"group_name"`
	Popularity    int    `db:"popularity"`
	AlbumImageURL string `db:"album_image_url"`
}

type Collection struct {
	Name        string    `db:"name"`
	GroupName   string    `db:"group_name"`
	ReleaseDate time.Time `db:"release_date"`
}

type CollectionWithCardCount struct {
	Collection
	CardCount int `db:"card_count"`
}

type Idol struct {
	StageName string    `db:"stage_name"`
	RealName  string    `db:"real_name"`
	GroupName string    `db:"group_name"`
	Birthday  time.Time `db:"birthday"`
	Country   string    `db:"country"`
	MBTI      string    `db:"mbti"`
}

type IdolWithBias struct {
	Idol
	Bias bool `db:"bias"`
}

type Card struct {
	ID             int    `db:"id"`
	Variant        string `db:"variant"`
	IdolName       string `db:"idol_name"`
	CollectionName string `db:"collection_name"`
	GroupName      string `db:"group_name"`
}

type CardWithOwnershipState struct {
	Card
	Owned bool `db:"owned"`
}

type User struct {
	Username  string `db:"username"`
	Email     string `db:"email"`
	Provider  string `db:"provider"`
	AvatarURL string `db:"avatar_url"`
	Bio       string `db:"bio"`
}

type UserWithFollowers struct {
	User
	IsFollowing bool `db:"is_following"`
	Followers   int  `db:"followers"`
	Following   int  `db:"following"`
}

type NewUser struct {
	Email     string `db:"email"`
	Provider  string `db:"provider"`
	AvatarURL string `db:"avatar_url"`
}
