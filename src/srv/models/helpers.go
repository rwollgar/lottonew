package models

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/src/srv/lottoerrors"
)

//func extractDataURL(games map[string]game, body io.ReadCloser) (map[string]game, error) {
func extractDataURL(games map[string]game, body io.ReadCloser) error {

	const op lottoerrors.Operation = "helpers:extractDataURL"

	doc, err := goquery.NewDocumentFromReader(body)

	if err != nil {
		return lottoerrors.E(err, lottoerrors.Severity4, lottoerrors.Category3, op, "Error reading document")
	}

	doc.Find(".lw-freqchart-list").Children().Each(func(i int, item *goquery.Selection) {

		href, exists := item.Find("a").Attr("href")

		if exists {
			href = strings.ToLower(href)
			for k, g := range games {
				if strings.Contains(href, g.Name) {
					fmt.Println(href)
					g.DataURL = href
					games[k] = g
				}
			}
		}

	})

	return nil
}

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

func generateBuckets(startDraw draw, numberOfDraws int, g game) *drawInfo {

	var di drawInfo

	di.Game = g.Name
	di.MaxNumber = g.MaxNumber
	di.MaxNumbers = g.StandardNumbers + g.MaxSupplementary
	di.DrawID = startDraw.DrawID

	buckets := make(map[int]draw)
	//unique := make(map[int]int)

	startIndex := startDraw.ID //len(g.draws) - startDraw
	endIndex := startIndex + numberOfDraws + 1

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

	//di.Metrics = di.getMetrics()
	di.uniqueMetrics()
	di.missingMetrics()
	//di.Metrics.Ranges = getNumberRanges(di.Unique, g.MaxNumber)

	return &di
}

func generateMetrics(numbers []float64, supplements int) *drawMetrics {

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

	return &m
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
