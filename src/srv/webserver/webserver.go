package webserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/browser"
	"github.com/rs/cors"
	//errors "github.com/src/srv/errors"
)

//ServerContext struct
type ServerContext struct {
	Router             *mux.Router
	RandomNumberAPIURL string
	RandomNumberAPIKEY string
	WebUI              bool
	WebServer          bool
	Port               int
	Cwd                string
	RootDir            string
}

//Webserver => webserver
func Webserver() {
	fmt.Println("Package: webserver")
}

//InitWebserver => Initialisze webserver and routes
func (s *ServerContext) InitWebserver() error {

	if s.WebUI || s.WebServer {

		s.Router.Handle("/api/games", handlers.LoggingHandler(os.Stdout, s.validateToken(s.handleGetGames()))).Methods("GET")
		s.Router.Handle("/api/games/{game}", handlers.LoggingHandler(os.Stdout, s.handleGetGame())).Methods("GET")
		s.Router.Handle("/api/games/{game}/draws", handlers.LoggingHandler(os.Stdout, s.handleGetDrawsForGame())).Methods("GET")
		s.Router.Handle("/api/games/{game}/draws/{drawid}", handlers.LoggingHandler(os.Stdout, s.handleGetDrawsForGame())).Methods("GET")
		s.Router.PathPrefix("/").Handler(http.StripPrefix("/web", http.FileServer(http.Dir("../web/spa"))))

		serverURL := fmt.Sprintf(":%d", s.Port)
		browserURL := fmt.Sprintf("http://localhost:%d/web", s.Port)

		if s.WebUI {
			err := browser.OpenURL(browserURL)
			log.Printf("Error opening browser: %s (%s)", browserURL, err)
		}

		_ = log.Output(0, fmt.Sprintf("\n\nRunning Web Server on :%d\nPress CTRL+C to stop.", s.Port))
		handler := cors.Default().Handler(s.Router)
		log.Fatal(http.ListenAndServe(serverURL, handler))

	}

	return nil

}

//func (s *ServerContext) validateToken(h http.HandlerFunc) http.HandlerFunc {
func (s *ServerContext) validateToken(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
		fmt.Println("Validating JWToken")
	}
}

//ReturnError => Return error from webserver
func ReturnError(w http.ResponseWriter, r *http.Request, status int, err error) {

	code := strings.Replace(fmt.Sprintf("%T", err), "main.", "", 1)

	msg := err.Error()

	w.WriteHeader(status)

	type resp = struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	_ = json.NewEncoder(w).Encode(struct {
		Errors []resp `json:"errors"`
	}{[]resp{
		{Code: code, Message: msg},
	}})

}
