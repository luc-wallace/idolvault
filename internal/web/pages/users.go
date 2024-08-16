package pages

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/luc-wallace/idolvault/internal/models"
	"github.com/luc-wallace/idolvault/internal/web/util"
)

var (
	userInfoTmpl      = loadTemplate("user-base.html", "user-info.html")
	userCardsTmpl     = loadTemplate("user-base.html", "user-cards.html")
	userFollowersTmpl = loadTemplate("user-base.html", "user-followers.html")
	userBiasesTmpl    = loadTemplate("user-base.html", "user-biases.html")
)

type userPageParams struct {
	util.PageParams
	User          *models.UserWithFollowers
	Self          bool
	InfoPage      bool
	CardsPage     bool
	FollowersPage bool
	BiasesPage    bool
}

func (p *pageProvider) GetUserInfoPage(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	params := util.GetParams(r)

	if ok, _ := util.ValidateUsername(username); !ok {
		w.WriteHeader(http.StatusNotFound)
		notFoundTmpl.Execute(w, params)
		return
	}

	var currentUser string
	if params.Session.LoggedIn {
		currentUser = params.Session.User.Username
	}

	ctx := context.Background()
	user, err := p.db.GetUserWithFollowing(ctx, username, currentUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		serverErrorTmpl.Execute(w, serverErrorParams{params, err.Error()})
		return
	} else if user == nil {
		w.WriteHeader(http.StatusNotFound)
		notFoundTmpl.Execute(w, params)
		return
	}
	if err := userInfoTmpl.Execute(w, userPageParams{
		PageParams: params,
		User:       user,
		Self:       params.Session.LoggedIn && strings.EqualFold(params.Session.User.Username, user.Username),
		InfoPage:   true,
	}); err != nil {
		log.Println(err)
	}
}

type userCardsPageParams struct {
	userPageParams
	Cards  []*models.Card
	Empty  bool
	Groups []string
}

func (p *pageProvider) GetUserCardsPage(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	params := util.GetParams(r)

	if ok, _ := util.ValidateUsername(username); !ok {
		w.WriteHeader(http.StatusNotFound)
		notFoundTmpl.Execute(w, params)
		return
	}

	var currentUser string
	if params.Session.LoggedIn {
		currentUser = params.Session.User.Username
	}

	ctx := context.Background()
	user, err := p.db.GetUserWithFollowing(ctx, username, currentUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		serverErrorTmpl.Execute(w, serverErrorParams{params, err.Error()})
		return
	} else if user == nil {
		w.WriteHeader(http.StatusNotFound)
		notFoundTmpl.Execute(w, params)
		return
	}

	cards, _ := p.db.GetUserCards(ctx, username)
	var empty bool
	if len(cards) < 1 {
		empty = true
	}

	var groups []string
Main:
	for _, card := range cards {
		for _, group := range groups {
			if card.GroupName == group {
				continue Main
			}
		}
		groups = append(groups, card.GroupName)
	}

	if err := userCardsTmpl.Execute(w, userCardsPageParams{
		userPageParams: userPageParams{PageParams: params,
			User:      user,
			Self:      params.Session.LoggedIn && strings.EqualFold(params.Session.User.Username, user.Username),
			CardsPage: true,
		},
		Cards:  cards,
		Empty:  empty,
		Groups: groups,
	}); err != nil {
		log.Println(err)
	}
}

type userFollowersPageParams struct {
	userPageParams
	Followers []*models.User
	Empty     bool
	Count     int
}

func (p *pageProvider) GetUserFollowersPage(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	params := util.GetParams(r)

	if ok, _ := util.ValidateUsername(username); !ok {
		w.WriteHeader(http.StatusNotFound)
		notFoundTmpl.Execute(w, params)
		return
	}

	var currentUser string
	if params.Session.LoggedIn {
		currentUser = params.Session.User.Username
	}

	ctx := context.Background()
	user, err := p.db.GetUserWithFollowing(ctx, username, currentUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		serverErrorTmpl.Execute(w, serverErrorParams{params, err.Error()})
		return
	} else if user == nil {
		w.WriteHeader(http.StatusNotFound)
		notFoundTmpl.Execute(w, params)
		return
	}

	followers, _ := p.db.GetUserFollowers(context.Background(), username)
	var empty bool
	if len(followers) < 1 {
		empty = true
	}

	if err := userFollowersTmpl.Execute(w, userFollowersPageParams{
		userPageParams: userPageParams{PageParams: params,
			User:          user,
			Self:          params.Session.LoggedIn && strings.EqualFold(params.Session.User.Username, user.Username),
			FollowersPage: true,
		},
		Followers: followers,
		Count:     len(followers),
		Empty:     empty,
	}); err != nil {
		log.Println(err)
	}
}

type userBiasesPageParams struct {
	userPageParams
	Biases []*models.Idol
	Empty  bool
}

func (p *pageProvider) GetUserBiasesPage(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	params := util.GetParams(r)

	if ok, _ := util.ValidateUsername(username); !ok {
		w.WriteHeader(http.StatusNotFound)
		notFoundTmpl.Execute(w, params)
		return
	}

	var currentUser string
	if params.Session.LoggedIn {
		currentUser = params.Session.User.Username
	}

	ctx := context.Background()
	user, err := p.db.GetUserWithFollowing(ctx, username, currentUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		serverErrorTmpl.Execute(w, serverErrorParams{params, err.Error()})
		return
	} else if user == nil {
		w.WriteHeader(http.StatusNotFound)
		notFoundTmpl.Execute(w, params)
		return
	}

	biases, _ := p.db.GetUserBiases(context.Background(), username)
	var empty bool
	if len(biases) < 1 {
		empty = true
	}
	if err := userBiasesTmpl.Execute(w, userBiasesPageParams{
		userPageParams: userPageParams{PageParams: params,
			User:       user,
			Self:       params.Session.LoggedIn && strings.EqualFold(params.Session.User.Username, user.Username),
			BiasesPage: true,
		},
		Biases: biases,
		Empty:  empty,
	}); err != nil {
		log.Println(err)
	}
}
