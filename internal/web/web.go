package web

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/luc-wallace/idolvault/internal/db"
	"github.com/luc-wallace/idolvault/internal/web/api"
	"github.com/luc-wallace/idolvault/internal/web/pages"
)

func New(conn *db.Conn, sess *scs.SessionManager) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Mount("/api", api.Mount(conn, sess))
	r.Mount("/", pages.Mount(conn, sess))

	dir := http.Dir("./internal/web/static")
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(dir)))

	return r
}
