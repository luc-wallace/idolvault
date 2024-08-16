package pages

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/luc-wallace/idolvault/internal/models"
	"github.com/luc-wallace/idolvault/internal/web/util"
)

var (
	cardPageTmpl          = loadTemplate("card.html")
	groupCardsContentTmpl = loadComponent("group-card-collection.html")
	userCardsContentTmpl  = loadComponent("user-card-collection.html")
)

type cardPageParams struct {
	util.PageParams
	Card *models.Card
}

func (p *pageProvider) GetCardPage(w http.ResponseWriter, r *http.Request) {
	params := util.GetParams(r)
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		notFoundTmpl.Execute(w, params)
		return
	}

	card, err := p.db.GetCard(context.Background(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		notFoundTmpl.Execute(w, params)
		return
	}

	if err := cardPageTmpl.Execute(w, cardPageParams{PageParams: params, Card: card}); err != nil {
		log.Println(err)
	}
}

type groupCardsComponentParams struct {
	util.PageParams
	Cards []*models.CardWithOwnershipState
	Empty bool
}

func (p *pageProvider) GetGroupCardsContent(w http.ResponseWriter, r *http.Request) {
	params := util.GetParams(r)
	groupName, _ := url.QueryUnescape(chi.URLParam(r, "group"))
	collectionName, _ := url.QueryUnescape(chi.URLParam(r, "collection"))
	query := r.URL.Query()
	idolName, _ := url.QueryUnescape(query.Get("idol"))
	owned := query.Get("owned")

	if r.Header.Get("hx-request") != "true" {
		w.WriteHeader(http.StatusNotFound)
		notFoundTmpl.Execute(w, params)
		return
	}

	var cards []*models.CardWithOwnershipState
	var err error
	var username string
	if params.Session.LoggedIn {
		username = params.Session.User.Username
	}

	if idolName == "" || idolName == "All" {
		cards, err = p.db.GetCollectionCardsWithOwnershipState(context.Background(), username, groupName, collectionName, "")
	} else {
		cards, err = p.db.GetCollectionCardsWithOwnershipState(context.Background(), username, groupName, collectionName, idolName)
	}

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "couldn't find collection")
		return
	}

	if owned == "true" || owned == "false" {
		o := owned == "true"

		c := []*models.CardWithOwnershipState{}
		for _, card := range cards {
			if card.Owned == o {
				c = append(c, card)
			}
		}
		cards = c
	}

	empty := false
	if len(cards) < 1 {
		empty = true
	}

	if err := groupCardsContentTmpl.Execute(w, groupCardsComponentParams{PageParams: params, Cards: cards, Empty: empty}); err != nil {
		log.Println(err)
	}
}

type userCardsComponentParams struct {
	Cards []*models.Card
	Empty bool
	Self  bool
}

func (p *pageProvider) GetUserCardsContent(w http.ResponseWriter, r *http.Request) {
	params := util.GetParams(r)
	username := chi.URLParam(r, "username")
	query := r.URL.Query()
	group, _ := url.QueryUnescape(query.Get("group"))
	collection, _ := url.QueryUnescape(query.Get("collection"))
	idol, _ := url.QueryUnescape(query.Get("idol"))

	if group == "All" {
		group = ""
	}
	if collection == "All" {
		collection = ""
	}
	if idol == "All" {
		idol = ""
	}

	cards, err := p.db.GetUserCardsWithFilter(context.Background(), username, group, collection, idol)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal server error: %sv", err)
		return
	}

	var empty bool
	if len(cards) < 1 {
		empty = true
	}

	if err := userCardsContentTmpl.Execute(w, userCardsComponentParams{
		Cards: cards,
		Empty: empty,
		Self:  params.Session.LoggedIn && strings.EqualFold(params.Session.User.Username, username),
	}); err != nil {
		log.Println(err)
	}
}
