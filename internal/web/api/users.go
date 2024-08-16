package api

import (
	"context"
	"fmt"
	"html"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/luc-wallace/idolvault/internal/web/util"
)

func (p *apiProvider) GetUsernameStatus(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	if username == "" {
		fmt.Fprintf(w, "")
		return
	}

	ok, msg := util.ValidateUsername(username)
	if !ok {
		fmt.Fprintf(w, `<p class="error">%s</p>`, msg)
		return
	}
	user, err := p.db.GetUserByUsername(context.Background(), username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `<p class="error">server error :(</p>`)
		return
	} else if user != nil {
		fmt.Fprint(w, `<p class="error">username taken :(</p>`)
		return
	}

	fmt.Fprintf(w, `<p>'%s' is available :)</p>`, html.EscapeString(username))
}

func (p *apiProvider) PostUserFollowers(w http.ResponseWriter, r *http.Request) {
	params := util.GetParams(r)
	following := chi.URLParam(r, "following")
	follower := params.Session.User.Username

	if following == follower {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "you cannot follow yourself")
		return
	}

	if err := p.db.CreateUserFollower(context.Background(), follower, following); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "bad request")
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `
	<button class="button unfollow" hx-delete="/api/users/%s/followers" hx-swap="outerHTML" hx-target="this">
		Unfollow
	</button>
`, following)
}

func (p *apiProvider) DeleteUserFollowers(w http.ResponseWriter, r *http.Request) {
	params := util.GetParams(r)
	following := chi.URLParam(r, "following")
	follower := params.Session.User.Username

	if following == follower {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "you cannot unfollow yourself")
		return
	}

	if err := p.db.DeleteUserFollower(context.Background(), follower, following); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "bad request")
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `
	<button class="button follow" hx-post="/api/users/%s/followers" hx-swap="outerHTML" hx-target="this">
		Follow
	</button>
`, following)
}

func (p *apiProvider) PostUserBias(w http.ResponseWriter, r *http.Request) {
	params := util.GetParams(r)
	group := chi.URLParam(r, "group")
	idol := chi.URLParam(r, "idol")
	username := params.Session.User.Username

	if err := p.db.CreateUserBias(context.Background(), username, idol, group); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "bad request")
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `
	<button class="button remove" hx-delete="/api/groups/%s/members/%s/bias" hx-swap="outerHTML" hx-target="this">
		<i class="fa fa-xmark"></i>
		Remove bias
	</button>`, group, idol)
}

func (p *apiProvider) DeleteUserBias(w http.ResponseWriter, r *http.Request) {
	params := util.GetParams(r)
	group := chi.URLParam(r, "group")
	idol := chi.URLParam(r, "idol")
	username := params.Session.User.Username

	if err := p.db.DeleteUserBias(context.Background(), username, idol, group); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "bad request")
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `
	<button class="button add" hx-post="/api/groups/%s/members/%s/bias" hx-swap="outerHTML" hx-target="this">
		<i class="fa fa-add"></i>
		Add bias
	</button>`, group, idol)
}
