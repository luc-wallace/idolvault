package pages

import (
	"net/http"

	"github.com/luc-wallace/idolvault/internal/web/util"
)

var homeTmpl = loadTemplate("index.html")

type homePageParams struct {
	util.PageParams
	CardIDs1 []int
	CardIDs2 []int
}

func (p *pageProvider) GetHomePage(w http.ResponseWriter, r *http.Request) {
	// var cardIDs []int
	// i := 1
	// for len(cardIDs) < 30 {
	// 	i += rand.IntN(25) + 1
	// 	cardIDs = append(cardIDs, i)
	// }
	cardIDs := []int{5, 9, 21, 37, 58, 80, 104, 111, 119, 144, 157, 159, 173, 183, 200, 216, 241, 246, 253, 264, 265, 268, 269, 272, 295, 306, 309, 318, 334, 341}
	homeTmpl.Execute(w, homePageParams{PageParams: util.GetParams(r), CardIDs1: cardIDs[:15], CardIDs2: cardIDs[15:]})
}
