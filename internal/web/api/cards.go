package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/luc-wallace/idolvault/internal/web/util"
)

func (p *apiProvider) PutCardsClaim(w http.ResponseWriter, r *http.Request) {
	params := util.GetParams(r)
	username := params.Session.User.Username
	w.Header().Add("Content-Type", "application/json")

	cardID, err := strconv.Atoi(r.FormValue("card_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error": "invalid card id"}`)
		return
	}

	ctx := context.Background()
	owned, err := p.db.GetCardOwnershipState(ctx, username, cardID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "internal server error: %v"}`, err)
		return
	}

	if owned {
		err = p.db.SetCardUnowned(ctx, username, cardID)
	} else {
		err = p.db.SetCardOwned(ctx, username, cardID)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "internal server error: %v"}`, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"new_state": %v}`, !owned)
}
