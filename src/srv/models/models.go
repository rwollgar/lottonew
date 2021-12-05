package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/pieterclaerhout/go-waitgroup"
)

//CmdArgs => command line arguments
type CmdArgs struct {
	Game         string
	GameType     string
	DrawOffset   int
	Draws        int
	UseWebUI     bool
	UseWebserver bool
	Port         int
	RapiKey      string
	RapiURL      string
	DataURL      string
}

//Game structure including list of draws
type game struct {
	Name             string       `json:"name"`
	DataURL          string       `json:"dataurl"`
	StandardNumbers  int          `json:"standardnumbers"`
	Supplementary    int          `json:"supplementary"`
	MaxNumber        int          `json:"maxnumber"`
	MaxSupplementary int          `json:"maxsupplementary"`
	Format           string       `json:"format"`
	DrawOffset       int          `json:"drawoffset"`
	NumDraws         int          `json:"numdraws"`
	Draws            map[int]draw `json:"draws,omitempty"`
	Order            int          `json:"order"`
	LastDraw         draw         `json:"lastdraw"`
}

//NumberRange => store number range
type numberRange struct {
	Name    string    `json:"name"`
	Start   int       `json:"start"`
	End     int       `json:"end"`
	Numbers []float64 `json:"numbers"`
}

//DrawMetrics struct
type drawMetrics struct {
	CrossSum     float64                `json:"crosssum"`
	Avg          float64                `json:"average"`
	CrossSumSupp float64                `json:"crosssumsupp"`
	AvgSupp      float64                `json:"averagesupp"`
	Ranges       map[string]numberRange `json:"ranges"`
}

//Draw struct
type draw struct {
	ID       int          `json:"id"`
	DrawID   int          `json:"drawid"`
	Date     time.Time    `json:"date"`
	YearDay  int          `json:"yearday"`
	Week     int          `json:"week"`
	Month    int          `josn:"month"`
	Quarter  int          `json:"quarter"`
	Year     int          `json:"year"`
	Numbers  []float64    `json:"numbers"`
	Metrics  *drawMetrics `json:"drawmetrics,omitempty"`
	DrawInfo *drawInfo    `json:"drawinfo,omitempty"`
	Game     *game        `json:"-"` //Don't export to JSON
}

//RandomRequest struct used by radmon number service
type randomRequest struct {
	Version string       `json:"jsonrpc"`
	Method  string       `json:"method"`
	ID      int          `json:"id"`
	Params  randomParams `json:"params"`
}

//RandomParams struct
type randomParams struct {
	APIKey      string `json:"apiKey"`
	N           int    `json:"n"`
	Min         int    `json:"min"`
	Max         int    `json:"max"`
	Replacement bool   `json:"replacement"`
}

//DrawInfo struct for saving draw details
type drawInfo struct {
	ID             int          `json:"id"`
	Game           string       `json:"game"`
	MaxNumber      int          `json:"maxnumber"`
	MaxNumbers     int          `json:"maxnumbers"`
	DrawID         int          `json:"drawid"`
	Buckets        map[int]draw `json:"buckets"`
	RBuckets       map[int]draw `json:"rbuckets"`
	Unique         []float64    `json:"unique"`
	Missing        []float64    `json:"missing"`
	UniqueMetrics  *drawMetrics `json:"uniquemetrics,omitempty"`
	MissingMetrics *drawMetrics `json:"missingmetrics,omitempty"`
}

var _games map[string]game

//var _draws map[int]draw

//InitGames => Initialise games from online data
func InitGames(dataDir string, url string, args CmdArgs) error {

	var e error
	//const op lottoerrors.Operation = "models:initgames"

	fmt.Println("Initialising Games...")
	e = nil

	_games = make(map[string]game)

	_games["oz-lotto"] = game{
		Name:             "oz-lotto",
		StandardNumbers:  7,
		Supplementary:    2,
		MaxNumber:        45,
		MaxSupplementary: 45,
		Format:           "csv",
		DrawOffset:       1,
		Order:            4}

	_games["powerball"] = game{
		Name:             "powerball",
		StandardNumbers:  7,
		Supplementary:    1,
		MaxNumber:        35,
		MaxSupplementary: 20,
		Format:           "csv",
		DrawOffset:       1,
		Order:            5}

	_games["saturday-lotto"] = game{
		Name:             "saturday-lotto",
		StandardNumbers:  6,
		Supplementary:    2,
		MaxNumber:        45,
		MaxSupplementary: 45,
		Format:           "csv",
		DrawOffset:       2,
		Order:            3}

	_games["monday-lotto"] = game{
		Name:             "monday-lotto",
		StandardNumbers:  6,
		Supplementary:    2,
		MaxNumber:        45,
		MaxSupplementary: 45,
		Format:           "csv",
		DrawOffset:       2,
		Order:            1}

	_games["wednesday-lotto"] = game{
		Name:             "wednesday-lotto",
		StandardNumbers:  6,
		Supplementary:    2,
		MaxNumber:        45,
		MaxSupplementary: 45,
		Format:           "csv",
		DrawOffset:       2,
		Order:            3}

	wg := waitgroup.NewWaitGroup(10)

	//Get the HTML from site
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	// document, err := goquery.NewDocumentFromReader(resp.Body)
	// if err != nil {
	// 	return lottoerrors.E(err, lottoerrors.Severity4, lottoerrors.Category3, op, "Error reading document")
	// }

	//Parse HTML and extract data urls
	//_games, err = extractDataURL(_games, resp.Body)
	err = extractDataURL(_games, resp.Body)

	if err != nil {
		return err
	}

	for ix, g := range _games {

		wg.BlockAdd()

		go func(g game, ix string) {

			defer wg.Done()

			resp, err := http.Get(g.DataURL)

			if err != nil {
				log.Fatal(err)
			}

			defer resp.Body.Close()

			dataErr := getData(&g, resp.Body)

			//Non fatal error. Continue with the data already downloaded.
			if dataErr != nil {
				fmt.Println(dataErr)
			}

			g.LastDraw = g.Draws[0]
			g.NumDraws = len(g.Draws)
			_games[ix] = g

		}(g, ix)

	}

	wg.Wait()
	fmt.Println("FINISHED")

	thisGame := _games[args.Game]
	thisDraw := thisGame.Draws[args.DrawOffset-1]
	thisDraw.DrawInfo = generateBuckets(thisDraw, args.Draws, thisGame)
	//generateBuckets(thisDraw, args.Draws, thisGame)
	//thisDraw.DrawInfo = &di
	_games[args.Game].Draws[args.DrawOffset-1] = thisDraw

	thisDraw.DrawInfo.printDrawInfo()
	fmt.Println()

	for _, g := range _games {
		g.printGame()
	}

	fmt.Printf("\n")

	thisDraw.DrawInfo.generateRandomPicks(args.RapiKey, args.RapiURL)

	return e

}

func GetGames() map[string]game {
	return _games
}

func GetGamesInfo() []game {

	i := 0
	temp := make([]game, len(_games))

	for _, g := range _games {

		temp[i] = g
		temp[i].Draws = map[int]draw{}
		i = i + 1
	}

	return temp

	//	return _games
}

func initDraw(d draw, record []string) draw {

	arlen := d.Game.StandardNumbers + d.Game.Supplementary
	d.Numbers = make([]float64, arlen)

	drawID, err := strconv.Atoi(record[0])

	if err != nil {
		d.DrawID = 0
	}

	d.DrawID = drawID
	d.Date, err = time.Parse("02/01/2006", record[1])
	yr, wk := d.Date.ISOWeek()
	m := d.Date.Month()
	d.YearDay = d.Date.YearDay()
	d.Week = wk
	d.Month = int(m)
	d.Quarter = int(math.Floor(float64(d.Month-1)/3)) + 1
	d.Year = yr

	if err != nil {
		fmt.Println(err)
	}

	var num float64

	offset := 2
	for i := 0; i < arlen; i++ {
		num, _ = strconv.ParseFloat(record[i+offset], 64)
		d.Numbers[i] = num
	}

	d.drawMetrics()
	d.Metrics.Ranges = getNumberRanges(d.Numbers, d.Game.MaxNumber)

	return d
}

func (di *drawInfo) uniqueMetrics() {
	di.UniqueMetrics = generateMetrics(di.Unique, 0)
	di.UniqueMetrics.Ranges = getNumberRanges(di.Unique, di.MaxNumber)
}

func (di *drawInfo) missingMetrics() {
	di.MissingMetrics = generateMetrics(di.Missing, 0)
	di.MissingMetrics.Ranges = getNumberRanges(di.Missing, di.MaxNumber)
}

func (d *draw) drawMetrics() {
	d.Metrics = generateMetrics(d.Numbers, d.Game.Supplementary)
}

func (di *drawInfo) generateRandomPicks(apikey string, url string) {

	if apikey != "" {

		r := &randomRequest{
			Method:  "generateIntegers",
			ID:      1,
			Version: "2.0",
			Params: randomParams{
				APIKey:      apikey,
				N:           7,
				Min:         1,
				Max:         45,
				Replacement: true}}

		b, err := json.Marshal(r)
		_ = err
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(b)))

		if err != nil {
			fmt.Println(err)
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		_ = err

		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))

	}

	rand.Seed(time.Now().UTC().UnixNano())

	n1 := rand.Intn(di.MaxNumber)
	n2 := rand.Intn(di.MaxNumber)
	n3 := rand.Intn(di.MaxNumber)

	fmt.Printf("Random Number: %d\n", n1)
	fmt.Printf("Random Number: %d\n", n2)
	fmt.Printf("Random Number: %d\n", n3)

}

func (g *game) GetDraw(drawID int) (draw, error) {

	var d draw

	for k := range g.Draws {
		if g.Draws[k].DrawID == drawID {
			return g.Draws[k], nil
		}
	}

	return d, fmt.Errorf("Can't find Draw with DarwID: %d", drawID)
}

func (g game) printGame() {
	fmt.Printf("%-15s\t%d\t%d\t%d\t%d\t%d\n", g.Name, g.StandardNumbers, g.Supplementary, g.MaxNumber, g.MaxSupplementary, len(g.Draws))
}

func (di drawInfo) printDrawInfo() {

	fmt.Printf("Game: %s. %d Draws\n\n", di.Game, di.DrawID)
	fmt.Printf("Numbers in set:\t\t%2v\n", di.Unique)
	fmt.Printf("Numbers missing:\t%2v\n\n", di.Missing)

	keys := make([]int, 0)
	for k := range di.Buckets {
		keys = append(keys, k)
	}
	//sort.Sort(sort.IntSlice(keys))
	sort.Ints(sort.IntSlice(keys))

	drawNumbers := make(map[int]int)
	for _, v := range di.Buckets[keys[0]].Numbers {
		drawNumbers[int(v)] = int(v)
	}

	for _, v := range keys {

		numbers := make([]float64, 0)
		numbers = append(numbers, di.Buckets[int(v)].Numbers...)
		fmt.Printf("Bucket %d [%d]\t\t%2v\n", di.Buckets[v].DrawID, len(numbers), numbers)
	}

	fmt.Printf("\nBucket Stats\n")
	// fmt.Printf("CrossSum:\t\t%d\n", int(di.Metrics.CrossSum))
	// fmt.Printf("Avg:\t\t\t%.2f\n\n", di.Metrics.Avg)

	// for _, v := range di.Metrics.Ranges {
	// 	fmt.Printf("%s\t%v\n", v.Name, v.Numbers)
	// }
}
