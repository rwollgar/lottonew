package models

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func ConfigureGames() map[string]game {

	games := make(map[string]game)

	games["oz-lotto"] = game{
		Name:             "oz-lotto",
		GameID:           "5130",
		DataURL:          "https://api.lotterywest.wa.gov.au/api/v1/games/5130/results-csv",
		StandardNumbers:  7,
		Supplementary:    2,
		MaxNumber:        45,
		MaxSupplementary: 45,
		Format:           "csv",
		DrawOffset:       1,
		Order:            4}

	games["powerball"] = game{
		Name:             "powerball",
		GameID:           "5132",
		DataURL:          "https://api.lotterywest.wa.gov.au/api/v1/games/5132/results-csv",
		StandardNumbers:  7,
		Supplementary:    1,
		MaxNumber:        35,
		MaxSupplementary: 20,
		Format:           "csv",
		DrawOffset:       1,
		Order:            5}

	games["saturday-lotto"] = game{
		Name:             "saturday-lotto",
		GameID:           "5127",
		DataURL:          "https://api.lotterywest.wa.gov.au/api/v1/games/5127/results-csv",
		StandardNumbers:  6,
		Supplementary:    2,
		MaxNumber:        45,
		MaxSupplementary: 45,
		Format:           "csv",
		DrawOffset:       2,
		Order:            3}

	games["monday-lotto"] = game{
		Name:             "monday-lotto",
		GameID:           "5128",
		DataURL:          "https://api.lotterywest.wa.gov.au/api/v1/games/5128/results-csv",
		StandardNumbers:  6,
		Supplementary:    2,
		MaxNumber:        45,
		MaxSupplementary: 45,
		Format:           "csv",
		DrawOffset:       2,
		Order:            1}

	games["wednesday-lotto"] = game{
		Name:             "wednesday-lotto",
		GameID:           "5129",
		DataURL:          "https://api.lotterywest.wa.gov.au/api/v1/games/5129/results-csv",
		StandardNumbers:  6,
		Supplementary:    2,
		MaxNumber:        45,
		MaxSupplementary: 45,
		Format:           "csv",
		DrawOffset:       2,
		Order:            3}

	return games
}

//func extractDataURL(games map[string]game, body io.ReadCloser) (map[string]game, error) {
// func setDataURL(games map[string]game, body io.ReadCloser) error {

// 	const op lottoerrors.Operation = "helpers:extractDataURL"

// 	doc, err := goquery.NewDocumentFromReader(body)

// 	if err != nil {
// 		return lottoerrors.E(err, lottoerrors.Severity4, lottoerrors.Category3, op, "Error reading document")
// 	}

// 	doc.Find(".lw-freqchart-list").Children().Each(func(i int, item *goquery.Selection) {

// 		href, exists := item.Find("a").Attr("href")

// 		if exists {
// 			href = strings.ToLower(href)
// 			for k, g := range games {
// 				if strings.Contains(href, g.GameID) {
// 					fmt.Println(href)
// 					g.DataURL = href
// 					games[k] = g
// 				}
// 			}
// 		}

// 	})

// 	return nil
// }

func getData(g *game, body io.ReadCloser) error {

	//fmt.Printf("Retrieving data for URL:%s", g.DataURL)
	//time.Sleep(time.Duration(30000) * time.Millisecond)

	// resp, err := http.Get(g.DataURL)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	//defer resp.Body.Close()

	reader := csv.NewReader(body)
	reader.FieldsPerRecord = -1

	g.Draws = make(map[int]draw)

	//Read header record
	_, err := reader.Read()
	if err == io.EOF {
		return nil
	}

	i := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err == nil {
			g.Draws[i] = initDraw(draw{Game: g, ID: i}, record)
			i++
		} else {
			log.Print(err)
		}

	}

	return err

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

	return d
}

// func generateBuckets(g game, startDraw draw, numberOfDraws int) *drawInfo {
func generateBuckets(startDraw *draw, numberOfDraws int) *drawInfo {

	var di drawInfo

	g := startDraw.Game //.Game = g.Name
	di.MaxNumber = g.MaxNumber
	di.MaxNumbers = g.StandardNumbers + g.MaxSupplementary
	di.DrawID = startDraw.DrawID

	buckets := make(map[int]draw)
	//unique := make(map[int]int)

	startIndex := startDraw.ID //len(g.draws) - startDraw
	endIndex := startIndex + numberOfDraws

	firstDraw := g.Draws[startIndex]

	var b draw
	//b.Numbers = make(map[int]int)
	b.Numbers = make([]float64, len(firstDraw.Numbers))
	unique := make(map[int]float64)

	b.DrawID = firstDraw.DrawID
	b.ID = 1
	b.Date = firstDraw.Date

	copy(b.Numbers, firstDraw.Numbers)

	for i := 0; i < len(firstDraw.Numbers); i++ {
		unique[int(firstDraw.Numbers[i])] = firstDraw.Numbers[i]
	}
	buckets[startIndex] = b

	for i := startIndex + 1; i < endIndex; i++ {

		d := g.Draws[i]
		var b draw

		b.Numbers = make([]float64, 0)
		b.DrawID = d.DrawID
		b.Date = d.Date
		b.ID = i

		//for _, n := range d.Numbers {
		for i := 0; i < len(d.Numbers); i++ {
			idx := int(d.Numbers[i])
			_, ok := unique[idx]

			if !ok {
				b.Numbers = append(b.Numbers, d.Numbers[i])
				unique[idx] = d.Numbers[i]
			}
		}

		buckets[i] = b

	}

	di.Buckets = buckets
	for _, v := range unique {
		di.Unique = append(di.Unique, v)
	}

	for i := 1; i <= g.MaxNumber; i++ {
		_, ok := unique[i]
		if !ok {
			di.Missing = append(di.Missing, float64(i))
		}
	}

	di.uniqueMetrics()
	di.missingMetrics()

	return &di
}

func generateMetrics(numbers []float64, supplements int) drawMetrics {

	var m drawMetrics

	sum := 0.0
	sumsupp := 0.0
	stdLength := len(numbers) - supplements

	for i := 0; i < len(numbers); i++ {

		sumsupp += numbers[i]
		if i < stdLength {
			sum += numbers[i]
		}
	}

	m.CrossSum = sum
	m.CrossSumSupp = sumsupp
	m.Avg = math.Round((float64(sum)/float64(stdLength))*100) / 100
	m.AvgSupp = math.Round((float64(sumsupp)/float64(len(numbers)))*100) / 100

	return m
}

func getNumberRanges(numbers []float64, maxNum int) map[string]numberRange {

	//Default Ranges 1TO9, 10TO19, 20TO29, 30TO39, 40TO45
	r := make(map[string]numberRange)

	offset := 9
	startNum := 1
	endNum := offset

	for i := 0; i < 5; i++ {

		nr := numberRange{
			Name:  fmt.Sprintf("%02dTO%02d", startNum, endNum),
			Start: startNum,
			End:   endNum,
		}

		//nr.Numbers = make(map[int]int)

		for i := 0; i < len(numbers); i++ {
			n := int(numbers[i])
			if n >= nr.Start && n <= nr.End {
				nr.Numbers = append(nr.Numbers, float64(n)) //[n] = n
			}
		}

		startNum = endNum + 1
		endNum = startNum + offset

		if endNum > maxNum {
			endNum = maxNum
		}

		r[nr.Name] = nr

		//fmt.Println(nr)
	}

	return r
}

func getRandomNumbers(d draw, apikey string, url string) []int {

	var randomNumbers = make([]int, d.Game.StandardNumbers+d.Game.Supplementary)

	_ = randomNumbers

	if apikey != "" {

		r := &randomRequest{
			Method:  "generateIntegers",
			ID:      1,
			Version: "2.0",
			Params: randomApiParams{
				APIKey:      apikey,
				N:           d.Game.StandardNumbers + d.Game.Supplementary,
				Min:         1,
				Max:         d.Game.MaxNumber,
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
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("\n\n\nresponse Body:", string(body))

		var data randomApiResponse
		if err := json.Unmarshal(body, &data); err != nil {
			fmt.Println("failed to unmarshal:", err)
		} else {
			fmt.Println(data)
			fmt.Printf("%s %s\n\n", d.Game.Name, data.Jsonrpc)
			fmt.Printf("%s %v", d.Game.Name, data.Result.Random.Data)

		}

	} else {

		rand.Seed(time.Now().UTC().UnixNano())

		n1 := rand.Intn(d.Game.MaxNumber)
		n2 := rand.Intn(d.Game.MaxNumber)
		n3 := rand.Intn(d.Game.MaxNumber)

		fmt.Printf("Random Number: %d\n", n1)
		fmt.Printf("Random Number: %d\n", n2)
		fmt.Printf("Random Number: %d\n", n3)

	}

	return randomNumbers
}
