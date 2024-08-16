package spotify

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/luc-wallace/idolvault/internal/db"
	"github.com/luc-wallace/idolvault/internal/models"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/clientcredentials"
)

type SpotifyService struct {
	client *spotify.Client
	db     *db.Conn
	config *clientcredentials.Config
}

func New(ctx context.Context, clientID string, clientSecret string, conn *db.Conn) (*SpotifyService, error) {
	s := &SpotifyService{db: conn, config: &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotifyauth.TokenURL,
	}}

	token, err := s.config.Token(ctx)
	httpClient := spotifyauth.New().Client(ctx, token)

	s.client = spotify.New(httpClient)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *SpotifyService) Start(ctx context.Context) {
	go s.refreshToken(ctx)
	for {
		groups, err := s.db.GetGroups(ctx)
		if err != nil {
			log.Fatal(err)
		}
		now := time.Now().UTC()
		for _, group := range groups {
			diff := now.Sub(group.UpdatedAt)
			if diff.Hours() < 1 {
				continue
			}
			if err := s.RefreshGroupStats(ctx, group); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("refreshed stats for %s\n", group.Name)
			time.Sleep(1 * time.Second)
		}
		time.Sleep(1 * time.Hour)
	}
}

func (s *SpotifyService) refreshToken(ctx context.Context) {
	for {
		time.Sleep(59 * time.Minute)
		token, err := s.config.Token(ctx)
		if err != nil {
			log.Printf("failed to fetch new spotify token: %v\n", token)
		}
		httpClient := spotifyauth.New().Client(ctx, token)
		s.client = spotify.New(httpClient)
		fmt.Println("refreshed spotify token")
	}
}

func (s *SpotifyService) RefreshGroupStats(ctx context.Context, group *models.Group) error {
	artistID := spotify.ID(group.SpotifyID)
	artist, err := s.client.GetArtist(ctx, artistID)
	if err != nil {
		return err
	}

	updated := &models.Group{
		Name:       group.Name,
		Popularity: int(artist.Popularity),
		Followers:  int64(artist.Followers.Count),
		ImageURL:   artist.Images[0].URL,
		Genres:     artist.Genres,
		UpdatedAt:  time.Now().UTC(),
	}
	if err := s.db.UpdateGroupStats(ctx, updated); err != nil {
		return err
	}

	tracks, err := s.client.GetArtistsTopTracks(ctx, artistID, "US")
	if err != nil {
		return err
	}

	if err := s.db.DeleteGroupSongs(ctx, group.Name); err != nil {
		return err
	}
	for _, track := range tracks {
		if err := s.db.CreateSong(ctx, &models.Song{
			SpotifyID:     track.ID.String(),
			Name:          track.Name,
			GroupName:     group.Name,
			Popularity:    int(track.Popularity),
			AlbumImageURL: track.Album.Images[0].URL,
		}); err != nil {
			return err
		}
	}
	return nil
}
