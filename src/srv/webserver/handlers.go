package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/src/srv/models"
)

//Handler1 => Handler
func Handler1() {
	fmt.Println("Package: handlers")
}

func (s ServerContext) handleGetGames() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// i := 0
		// g := make([]game, len(models.GetGames()))

		// for _, v := range models.GetGames() {
		// 	g[i] = v
		// 	i++
		// }

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(models.GetGamesInfo())

		fmt.Println("handleGetGames")

	}
}

func (s ServerContext) handleGetGame() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		g, ok := models.GetGames()[params["game"]]

		w.Header().Set("Content-Type", "application/json")

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error: Game <%s> not found.", params["game"]))
			return
		}

		_ = json.NewEncoder(w).Encode(g)
		fmt.Println("handleGetGame")

	}
}

func (s ServerContext) handleGetDrawsForGame() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		g, okGame := models.GetGames()[params["game"]]

		w.Header().Set("Content-Type", "application/json")

		if !okGame {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error: Draws for game <%s> not found.", params["game"]))
			return
		}

		drawid, okDrawid := strconv.Atoi(params["drawid"])

		if okDrawid != nil {

			_ = json.NewEncoder(w).Encode(g.Draws)
			return
		}

		for _, d := range g.Draws {

			if d.DrawID == drawid {

				_ = json.NewEncoder(w).Encode(d)
				return

			}
		}

		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(fmt.Sprintf("Error: Draw for game <%s> and draw Id <%s> not found.", params["game"], params["drawid"]))

		fmt.Println("handleGetDrawsForGame")

	}

}
