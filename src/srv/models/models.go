package models

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
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
	StaticDir    string
	JwtKey       string
}

//Game structure including list of draws
type game struct {
	Name             string       `json:"name"`
	DataURL          string       `json:"dataurl,omitempty"`
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

//Draw struct
type draw struct {
	ID      int         `json:"id"`
	DrawID  int         `json:"drawid"`
	Date    time.Time   `json:"date"`
	YearDay int         `json:"yearday"`
	Week    int         `json:"week"`
	Month   int         `josn:"month"`
	Quarter int         `json:"quarter"`
	Year    int         `json:"year"`
	Numbers []float64   `json:"numbers"`
	Metrics drawMetrics `json:"drawmetrics,omitempty"`
	SetInfo *drawInfo   `json:"drawinfo,omitempty"`
	Game    *game       `json:"-"` //Don't export to JSON
}

//DrawInfo struct for saving draw details
type drawInfo struct {
	ID             int          `json:"id"`
	MaxNumber      int          `json:"maxnumber"`
	MaxNumbers     int          `json:"maxnumbers"`
	DrawOffset     int          `json:"drawoffset"`
	NumDraws       int          `json:"numdraws"`
	DrawID         int          `json:"drawid"`
	Buckets        map[int]draw `json:"buckets"`
	RBuckets       map[int]draw `json:"rbuckets"`
	Unique         []float64    `json:"unique"`
	Missing        []float64    `json:"missing"`
	UniqueMetrics  drawMetrics  `json:"uniquemetrics,omitempty"`
	MissingMetrics drawMetrics  `json:"missingmetrics,omitempty"`
}

//DrawMetrics struct
type drawMetrics struct {
	CrossSum     float64                `json:"crosssum"`
	Avg          float64                `json:"average"`
	CrossSumSupp float64                `json:"crosssumsupp"`
	AvgSupp      float64                `json:"averagesupp"`
	Ranges       map[string]numberRange `json:"ranges"`
}

//RandomRequest struct used by radmon number service
type randomRequest struct {
	Version string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	ID      int             `json:"id"`
	Params  randomApiParams `json:"params"`
}

//RandomParams struct
type randomApiParams struct {
	APIKey      string `json:"apiKey"`
	N           int    `json:"n"`
	Min         int    `json:"min"`
	Max         int    `json:"max"`
	Replacement bool   `json:"replacement"`
}

type randomApiResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Random struct {
			Data           []int  `json:"data"`
			CompletionTime string `json:"completionTime"`
		} `json:"random"`
		BitsUsed      int `json:"bitsUsed"`
		BitsLeft      int `json:"bitsLeft"`
		RequestsLeft  int `json:"requestsLeft"`
		AdvisoryDelay int `json:"advisoryDelay"`
	} `json:"result"`
	ID int `json:"id"`
}

var _games map[string]game

//InitGames => Initialise games from online data
func InitGames(dataDir string, args CmdArgs) error {

	var e error
	//const op lottoerrors.Operation = "models:initgames"

	fmt.Println("Initialising Game data...")
	e = nil

	wg := waitgroup.NewWaitGroup(10)

	//Get the HTML from site
	resp, err := http.Get(args.DataURL)

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
	err = setDataURL(_games, resp.Body)

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
	thisDraw.generateSetInfo(args.Draws) // = generateBuckets(thisGame, thisDraw, args.Draws)
	//generateBuckets(thisDraw, args.Draws, thisGame)
	//thisDraw.DrawInfo = &di
	_games[args.Game].Draws[args.DrawOffset-1] = thisDraw

	thisDraw.printDraw()
	fmt.Println()

	for _, g := range _games {
		g.printGame()
	}

	fmt.Printf("\n")

	thisDraw.generateRandomPicks(args.RapiKey, args.RapiURL)

	return e

}

func GetGames() map[string]game {
	return _games
}

func SetGames(games map[string]game) {
	_games = games
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

func (di *drawInfo) uniqueMetrics() {
	di.UniqueMetrics = generateMetrics(di.Unique, 0)
	di.UniqueMetrics.Ranges = getNumberRanges(di.Unique, di.MaxNumber)
}

func (di *drawInfo) missingMetrics() {
	di.MissingMetrics = generateMetrics(di.Missing, 0)
	di.MissingMetrics.Ranges = getNumberRanges(di.Missing, di.MaxNumber)
}

func (d *draw) generateSetInfo(setSize int) {

	d.SetInfo = generateBuckets(d, setSize)
}

func (d *draw) drawMetrics() {
	d.Metrics = generateMetrics(d.Numbers, d.Game.Supplementary)
	d.Metrics.Ranges = getNumberRanges(d.Numbers, d.Game.MaxNumber)

}

func (d *draw) generateRandomPicks(apikey string, url string) {

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

func (d draw) printDraw() {

	di := d.SetInfo

	fmt.Printf("Game: %s. %d Draws\n\n", strings.ToUpper(d.Game.Name), di.DrawID)
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

// func initDraw(d draw, record []string) draw {

// 	arlen := d.Game.StandardNumbers + d.Game.Supplementary
// 	d.Numbers = make([]float64, arlen)

// 	drawID, err := strconv.Atoi(record[0])

// 	if err != nil {
// 		d.DrawID = 0
// 	}

// 	d.DrawID = drawID
// 	d.Date, err = time.Parse("02/01/2006", record[1])
// 	yr, wk := d.Date.ISOWeek()
// 	m := d.Date.Month()
// 	d.YearDay = d.Date.YearDay()
// 	d.Week = wk
// 	d.Month = int(m)
// 	d.Quarter = int(math.Floor(float64(d.Month-1)/3)) + 1
// 	d.Year = yr

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	var num float64

// 	offset := 2
// 	for i := 0; i < arlen; i++ {
// 		num, _ = strconv.ParseFloat(record[i+offset], 64)
// 		d.Numbers[i] = num
// 	}

// 	d.drawMetrics()

// 	return d
// }
