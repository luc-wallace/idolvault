package pages

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/luc-wallace/idolvault/internal/models"
	"github.com/luc-wallace/idolvault/internal/web/util"
	"github.com/markbates/goth/gothic"
)

var (
	loginPageTmpl    = loadTemplate("login.html")
	registerPageTmpl = loadTemplate("register.html")
)

type loginPageParams struct {
	util.PageParams
	Error   bool
	Message string
}

func (p *pageProvider) GetLoginPage(w http.ResponseWriter, r *http.Request) {
	params := loginPageParams{PageParams: util.GetParams(r)}
	params.Message = p.sess.PopString(r.Context(), "error")
	if params.Message != "" {
		params.Error = true
	}

	loginPageTmpl.Execute(w, params)
}

func (p *pageProvider) Logout(w http.ResponseWriter, r *http.Request) {
	p.sess.Destroy(r.Context())
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (p *pageProvider) GetProvider(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
}

func WriteAvatar(imageURL string, username string) error {
	res, err := http.Get(imageURL)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer res.Body.Close()

	file, err := os.Create(fmt.Sprintf("./internal/web/static/img/users/%s.png", username))
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return fmt.Errorf("could not copy data: %v", err)
	}

	return nil
}

func (p *pageProvider) GetProviderCallback(w http.ResponseWriter, r *http.Request) {
	u, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	ctx := context.Background()

	// p.sess.Remove(ctx, "user")
	user, err := p.db.GetUserByEmail(ctx, u.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		serverErrorTmpl.Execute(nil, serverErrorParams{util.GetParams(r), err.Error()})
		return
	} else if user != nil { // User with email exists in the database
		// Email matches, but provider does not
		if user.Provider != u.Provider {
			p.sess.Put(r.Context(), "error", "Email already registered with a different provider.")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if err := WriteAvatar(u.AvatarURL, user.Username); err != nil {
			p.sess.Put(r.Context(), "error", fmt.Sprintf("Failed to refresh profile picture: %v", err))
		}
		p.sess.Put(r.Context(), "user", user)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	newUser := &models.NewUser{
		Email:     u.Email,
		Provider:  u.Provider,
		AvatarURL: u.AvatarURL,
	}

	p.sess.Put(r.Context(), "newUser", newUser)
	http.Redirect(w, r, "/register", http.StatusSeeOther)
}

func (p *pageProvider) GetRegisterPage(w http.ResponseWriter, r *http.Request) {
	registerPageTmpl.Execute(w, util.GetParams(r))
}

func (p *pageProvider) PostUsers(w http.ResponseWriter, r *http.Request) {
	if !p.sess.Exists(r.Context(), "newUser") {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	newUser := p.sess.Get(r.Context(), "newUser").(models.NewUser)
	username := r.FormValue("username")
	bio := r.FormValue("bio")

	ok, _ := util.ValidateUsername(username)
	if !ok {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}
	u, err := p.db.GetUserByUsername(context.Background(), username)
	if err != nil || u != nil {
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	p.sess.Remove(r.Context(), "newUser")

	user := &models.User{
		Username:  username,
		Email:     newUser.Email,
		Provider:  newUser.Provider,
		AvatarURL: fmt.Sprintf("/static/img/users/%s.png", username),
		Bio:       bio,
	}

	if err := p.db.CreateUser(context.Background(), user); err == nil {
		p.sess.Put(r.Context(), "user", user)
		WriteAvatar(newUser.AvatarURL, username)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
