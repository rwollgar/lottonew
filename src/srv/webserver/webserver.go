package webserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/browser"

	"github.com/src/srv/models"

	//errors "github.com/src/srv/errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//ServerContext struct
type ServerContext struct {
	Args               models.CmdArgs
	RandomNumberAPIURL string
	Cwd                string
	RootDir            string
}

//Webserver => webserver
func Webserver() {
	fmt.Println("Package: webserver")
}

//InitWebserver => Initialisze webserver and routes
func (s *ServerContext) InitWebserver() error {

	if s.Args.UseWebUI || s.Args.UseWebserver {

		e := echo.New()

		// s.Context = echo.New().NewContext()
		e.Use(middleware.CORS())
		e.Use(middleware.Logger())

		e.Static("assets", s.Args.StaticDir)

		g := e.Group("api/")

		g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
			Skipper:     jwtSkipper,
			TokenLookup: "query:token",
			SigningKey:  s.Args.JwtKey,
		}))

		g.GET("games", s.getGames)
		g.GET("games/:game", s.getGames)
		g.GET("games/:game/draws", s.getDrawsForGame)
		g.GET("games/:game/draws/:draw", s.getDrawsForGame)

		g.GET("metrics/:game/:draw/:draws", s.getMetrics)
		g.GET("metrics/:game/:draw", s.getMetrics)
		g.GET("metrics/:game", s.getMetrics)

		serverURL := fmt.Sprintf("localhost:%d", s.Args.Port)
		browserURL := fmt.Sprintf("http://localhost:%d/web", s.Args.Port)

		if s.Args.UseWebUI {
			err := browser.OpenURL(browserURL)
			log.Printf("Error opening browser: %s (%s)", browserURL, err)
		}

		_ = log.Output(0, fmt.Sprintf("\n\nRunning Web Server on http://%s\nPress CTRL+C to stop.", serverURL))

		e.Logger.Fatal(e.Start(serverURL))

	}

	return nil

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

func jwtSkipper(c echo.Context) bool {

	return true

}
