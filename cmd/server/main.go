package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/luc-wallace/idolvault/internal/db"
	"github.com/luc-wallace/idolvault/internal/spotify"
	"github.com/luc-wallace/idolvault/internal/web"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/google"
	spotifyauth "github.com/markbates/goth/providers/spotify"
)

var dev bool

func init() {
	flag.BoolVar(&dev, "dev", false, "developer mode")
}

func main() {
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file\n")
	}
	fmt.Println("loaded .env")

	domain := os.Getenv("DOMAIN")
	port := os.Getenv("PORT")

	if dev {
		fmt.Println("starting app in developer mode, omit the -dev flag for production mode")
		domain = "localhost:" + port
	} else {
		fmt.Println("starting app in production mode, use the -dev flag for developer mode")
		if domain == "" {
			log.Fatal(".env file missing DOMAIN value")
		}
	}

	ctx := context.Background()
	conn, err := db.Connect(ctx, os.Getenv("POSTGRES_URI"))
	if err != nil {
		log.Fatal("error connecting to database: " + err.Error() + "\n")
	}
	defer conn.Close(ctx)
	fmt.Println("connected to postgres")

	service, err := spotify.New(ctx, os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"), conn)
	if err != nil {
		log.Fatal("error connecting to spotify: " + err.Error() + "\n")
	}
	fmt.Println("connected to spotify")
	go service.Start(ctx)
	fmt.Println("started spotify service")

	cookieStore := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	cookieStore.MaxAge(60 * 60 * 24)
	cookieStore.Options.Path = "/"
	cookieStore.Options.Secure = !dev
	cookieStore.Options.HttpOnly = true
	gothic.Store = cookieStore

	var baseCallbackURI string
	if dev {
		baseCallbackURI = "http://" + domain + "/auth/callback?provider="
	} else {
		baseCallbackURI = "https://" + domain + "/auth/callback?provider="
	}

	goth.UseProviders(
		google.New(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			baseCallbackURI+"google",
		),
		discord.New(
			os.Getenv("DISCORD_CLIENT_ID"),
			os.Getenv("DISCORD_CLIENT_SECRET"),
			baseCallbackURI+"discord",
			discord.ScopeIdentify, discord.ScopeEmail,
		),
		spotifyauth.New(
			os.Getenv("SPOTIFY_CLIENT_ID"),
			os.Getenv("SPOTIFY_CLIENT_SECRET"),
			baseCallbackURI+"spotify",
		),
	)

	sessionManager := scs.New()
	sessionManager.Lifetime = 2 * 24 * time.Hour
	sessionManager.Store = pgxstore.New(conn.Pool())
	sessionManager.Cookie.Secure = !dev

	fmt.Println("running on port " + port)
	r := web.New(conn, sessionManager)

	if err := http.ListenAndServe(":"+port, sessionManager.LoadAndSave(r)); err != nil {
		log.Fatal(err)
	}
}
