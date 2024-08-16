package pages

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/luc-wallace/idolvault/internal/models"
	"github.com/luc-wallace/idolvault/internal/web/util"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	topGroupsTmpl    = loadTemplate("top-groups.html")
	groupInfoTmpl    = loadTemplate("group-base.html", "group-info.html")
	groupMusicTmpl   = loadTemplate("group-base.html", "group-music.html")
	groupCardsTmpl   = loadTemplate("group-base.html", "group-cards.html")
	idolTmpl         = loadTemplate("idol.html")
	groupFiltersTmpl = loadComponent("group-filters.html")
)

type topGroupsPageParams struct {
	util.PageParams
	Groups []*models.Group
}

func (p *pageProvider) GetTopGroupsPage(w http.ResponseWriter, r *http.Request) {
	params := util.GetParams(r)
	groups, err := p.db.GetGroups(context.Background())
	if err != nil {
		log.Println(err)
		serverErrorTmpl.Execute(w, serverErrorParams{params, err.Error()})
		return
	}

	topGroupsTmpl.Execute(w, topGroupsPageParams{PageParams: params, Groups: groups})
}

type groupPageParams struct {
	util.PageParams
	Group     *models.Group
	GroupID   string
	InfoPage  bool
	MusicPage bool
	CardsPage bool
}

type groupInfoPageParams struct {
	groupPageParams
	Members      []*models.Idol
	MembersCount int
}

func (p *pageProvider) GetGroupInfoPage(w http.ResponseWriter, r *http.Request) {
	groupName := chi.URLParam(r, "id")
	ctx := context.Background()
	params := util.GetParams(r)

	group, err := p.db.GetGroup(ctx, groupName)
	if err != nil || group == nil {
		w.WriteHeader(http.StatusNotFound)
		notFoundTmpl.Execute(w, params)
		return
	}

	members, err := p.db.GetGroupIdols(ctx, groupName)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		serverErrorTmpl.Execute(w, serverErrorParams{params, err.Error()})
		return
	}

	if err := groupInfoTmpl.Execute(w, groupInfoPageParams{
		groupPageParams: groupPageParams{
			PageParams: params,
			Group:      group,
			InfoPage:   true,
			GroupID:    groupName,
		},
		Members:      members,
		MembersCount: len(members),
	}); err != nil {
		log.Println(err)
	}
}

type groupMusicPageParams struct {
	groupPageParams
	Genres    string
	Followers string
	TopSongs  []*models.Song
	UpdatedAt string
}

func (p *pageProvider) GetGroupMusicPage(w http.ResponseWriter, r *http.Request) {
	groupName := chi.URLParam(r, "id")
	ctx := context.Background()
	params := util.GetParams(r)

	group, err := p.db.GetGroup(ctx, groupName)
	if err != nil || group == nil {
		w.WriteHeader(http.StatusNotFound)
		notFoundTmpl.Execute(w, params)
		return
	}

	songs, err := p.db.GetGroupSongs(ctx, groupName)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		serverErrorTmpl.Execute(w, serverErrorParams{params, err.Error()})
		return
	}

	printer := message.NewPrinter(language.English)

	if err := groupMusicTmpl.Execute(w, groupMusicPageParams{
		groupPageParams: groupPageParams{
			PageParams: params,
			Group:      group,
			MusicPage:  true,
			GroupID:    groupName,
		},
		Followers: printer.Sprintf("%d\n", group.Followers),
		TopSongs:  songs,
		Genres:    strings.Join(group.Genres, ", "),
		UpdatedAt: group.UpdatedAt.UTC().Format("2006-01-02 15:04"),
	}); err != nil {
		log.Println(err)
	}
}

type groupCardsPageParams struct {
	groupPageParams
	Collections []*models.CollectionWithCardCount
	Members     []*models.Idol
}

func (p *pageProvider) GetGroupCardsPage(w http.ResponseWriter, r *http.Request) {
	groupName := chi.URLParam(r, "id")
	ctx := context.Background()
	params := util.GetParams(r)

	group, err := p.db.GetGroup(ctx, groupName)
	if err != nil || group == nil {
		w.WriteHeader(http.StatusNotFound)
		notFoundTmpl.Execute(w, params)
		return
	}

	collections, err := p.db.GetGroupCollectionsWithCardCounts(ctx, groupName, true)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		serverErrorTmpl.Execute(w, serverErrorParams{params, err.Error()})
		return
	}

	members, _ := p.db.GetGroupIdols(ctx, groupName)

	if err := groupCardsTmpl.Execute(w, groupCardsPageParams{
		groupPageParams: groupPageParams{
			PageParams: params,
			Group:      group,
			CardsPage:  true,
			GroupID:    groupName,
		},
		Collections: collections,
		Members:     members,
	}); err != nil {
		log.Println(err)
	}
}

type idolPageParams struct {
	util.PageParams
	Idol     *models.IdolWithBias
	Birthday string
}

func (p *pageProvider) GetGroupMemberPage(w http.ResponseWriter, r *http.Request) {
	groupName := chi.URLParam(r, "group")
	idolName := chi.URLParam(r, "idol")
	params := util.GetParams(r)

	var username string
	if params.Session.LoggedIn {
		username = params.Session.User.Username
	}

	idol, err := p.db.GetGroupIdolWithBias(context.Background(), groupName, idolName, username)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		notFoundTmpl.Execute(w, params)
		return
	}

	getOrdinal := func(day int) string {
		if day >= 11 && day <= 13 {
			return "th"
		}
		switch day % 10 {
		case 1:
			return "st"
		case 2:
			return "nd"
		case 3:
			return "rd"
		default:
			return "th"
		}
	}
	birthday := fmt.Sprintf("%d%s %s %d", idol.Birthday.Day(), getOrdinal(idol.Birthday.Day()), idol.Birthday.Month().String(), idol.Birthday.Year())

	if err := idolTmpl.Execute(w, idolPageParams{PageParams: params, Idol: idol, Birthday: birthday}); err != nil {
		log.Println(err)
	}
}

type groupFiltersParams struct {
	Members     []*models.Idol
	Collections []*models.Collection
}

func (p *pageProvider) GetGroupFilters(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	idols, _ := p.db.GetGroupIdols(context.Background(), group)
	if len(idols) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "group not found")
		return
	}

	collections, _ := p.db.GetGroupCollections(context.Background(), group, true)
	groupFiltersTmpl.Execute(w, groupFiltersParams{Members: idols, Collections: collections})
}
