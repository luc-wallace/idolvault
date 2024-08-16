package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/luc-wallace/idolvault/internal/db"
	"github.com/luc-wallace/idolvault/internal/web/util"
)

type apiProvider struct {
	db   *db.Conn
	sess *scs.SessionManager
}

func Mount(conn *db.Conn, sess *scs.SessionManager) *chi.Mux {
	p := &apiProvider{conn, sess}
	r := chi.NewRouter()

	r.Use(util.UseDefaultParams(sess))
	r.Use(util.HTMXOnly)

	r.Get("/users/usernamestatus", p.GetUsernameStatus)

	// Only if authorized
	r.Group(func(r chi.Router) {
		r.Use(util.CheckAuthAPI(sess, true))

		r.Put("/cards/claim", p.PutCardsClaim)
		r.Post("/users/{following}/followers", p.PostUserFollowers)
		r.Delete("/users/{following}/followers", p.DeleteUserFollowers)
		r.Post("/groups/{group}/members/{idol}/bias", p.PostUserBias)
		r.Delete("/groups/{group}/members/{idol}/bias", p.DeleteUserBias)
	})

	return r
}
