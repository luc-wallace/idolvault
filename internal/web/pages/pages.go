package pages

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/luc-wallace/idolvault/internal/db"
	"github.com/luc-wallace/idolvault/internal/web/util"
)

var (
	notFoundTmpl    = loadTemplate("not-found.html")
	serverErrorTmpl = loadTemplate("server-error.html")
)

type pageProvider struct {
	db   *db.Conn
	sess *scs.SessionManager
}

type serverErrorParams struct {
	util.PageParams
	Error string
}

func Mount(conn *db.Conn, sess *scs.SessionManager) *chi.Mux {
	r := chi.NewRouter()
	p := pageProvider{conn, sess}

	r.Use(util.UseDefaultParams(sess))
	r.Get("/logout", p.Logout)

	r.Group(func(r chi.Router) {
		r.Use(util.RegisterRedirect(sess))
		r.Get("/", p.GetHomePage)

		// Unauthenticated only
		r.Group(func(r chi.Router) {
			r.Use(util.CheckAuth(sess, false))
			r.Post("/auth/users", p.PostUsers)
			r.Group(func(r chi.Router) {
				r.Use(util.RegisterRedirect(sess))
				r.Get("/login", p.GetLoginPage)
				r.Get("/register", p.GetRegisterPage)
				r.Get("/auth", p.GetProvider)
				r.Get("/auth/callback", p.GetProviderCallback)
			})
		})

		// Groups
		r.Get("/groups", p.GetTopGroupsPage)
		r.Get("/groups/{id}", p.GetGroupInfoPage)
		r.Get("/groups/{id}/music", p.GetGroupMusicPage)
		r.Get("/groups/{id}/cards", p.GetGroupCardsPage)
		r.Get("/groups/{group}/members/{idol}", p.GetGroupMemberPage)
		r.Get("/content/groupfilters", p.GetGroupFilters)

		// Cards
		r.Get("/content/groups/{group}/collections/{collection}/cards", p.GetGroupCardsContent)
		r.Get("/content/users/{username}/cards", p.GetUserCardsContent)
		r.Get("/cards/{id}", p.GetCardPage)

		// Users
		r.Get("/users/{username}", p.GetUserInfoPage)
		r.Get("/users/{username}/cards", p.GetUserCardsPage)
		r.Get("/users/{username}/bias", p.GetUserBiasesPage)
		r.Get("/users/{username}/followers", p.GetUserFollowersPage)

	})

	r.NotFound(p.NotFoundHandler)

	return r
}

func loadTemplate(files ...string) *template.Template {
	paths := []string{"./internal/web/pages/templates/base.html"}
	for _, file := range files {
		paths = append(paths, fmt.Sprintf("./internal/web/pages/templates/%s", file))
	}
	return template.Must(template.ParseFiles(paths...))
}

func loadComponent(file string) *template.Template {
	return template.Must(template.ParseFiles(fmt.Sprintf("./internal/web/pages/content/%s", file)))
}

func (p *pageProvider) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	notFoundTmpl.Execute(w, util.GetParams(r))
}
