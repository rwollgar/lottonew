package webserver

import (
	//"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/src/srv/models"
)

func (s ServerContext) getGames(c echo.Context) error {

	p := c.Param("game")

	if p == "" {
		return c.JSON(http.StatusOK, models.GetGamesInfo())
	}

	g, ok := models.GetGames()[c.Param("game")]

	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Error: Game <%s> not found.", c.Param("game")))
	}

	return c.JSON(http.StatusOK, g)

}

func (s ServerContext) getDrawsForGame(c echo.Context) error {

	g, ok := models.GetGames()[c.Param("game")]

	if !ok {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Error: Game <%s> not found.", c.Param("game")))
	}

	drawid, err := strconv.Atoi(c.Param("draw"))

	if err != nil {

		return c.JSON(http.StatusOK, g.Draws)

	}

	for _, d := range g.Draws {

		if d.DrawID == drawid {

			return c.JSON(http.StatusOK, d)

		}
	}

	return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Error: Draw for game <%s> and draw Id <%s> not found.", c.Param("game"), c.Param("draw")))

}

func (s ServerContext) getMetrics(c echo.Context) error {

	g := c.Param("game")
	_ = g

	draw, err := strconv.Atoi(c.Param("draw"))
	_ = draw

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Error: Invalid draw <%s> not found.", c.Param("draw")))
	}

	fmt.Println("handleGetDrawsForGame")

	return c.String(http.StatusOK, "OK")

}
