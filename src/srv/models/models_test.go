package models

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
)

func TestModels(t *testing.T) {

	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Domain Models Test Suite")

	// num := testMe(99)

	// if num != 99 {
	// 	t.Errorf("Number, got: %d, want %d", num, 99)
	// }
}

var _ = ginkgo.Describe("Domain Models", func() {

	ginkgo.Context("initialise", func() {

		file, err := os.Open("../fixtures/ozlotto-30-11-2021-18_33.csv")
		fmt.Println(err)
		rdr := bufio.NewReader(file)
		stringReadCloser := io.NopCloser(rdr)

		ozLotto := game{
			Name:             "oz-lotto",
			StandardNumbers:  7,
			Supplementary:    2,
			MaxNumber:        45,
			MaxSupplementary: 45,
			Format:           "csv",
			DrawOffset:       1,
			Order:            4}

		dataErr := getData(&ozLotto, stringReadCloser)
		fmt.Println(dataErr)

		// record := []string{
		// 	"1350", "31/12/2019", "10", "11", "29", "1", "39", "23", "5",
		// 	"42", "40", "0", "0", "0.00", "0.00", "1", "0", "123427.60",
		// 	"123427.60", "66", "11", "254116.50", "3850.25", "343", "35",
		// 	"130683.00", "381.00", "3864", "461", "152434.80", "39.45",
		// 	"79375", "9827", "1742281.25", "21.95", "118330", "14958", "1952445.00",
		// 	"16.50"}

		// g := game{
		// 	Name:             "oz-lotto",
		// 	StandardNumbers:  7,
		// 	Supplementary:    2,
		// 	MaxNumber:        45,
		// 	MaxSupplementary: 45,
		// 	Format:           "csv",
		// 	DrawOffset:       1}

		ginkgo.It("should return 5 games", func() {
			games := ConfigureGames()
			gomega.Expect(len(games)).To(gomega.Equal(5))
		})

		ginkgo.It("should return draw 1350", func() {
			d, e := ozLotto.GetDraw(1350)
			gomega.Expect(d).NotTo(gomega.BeNil())
			gomega.Expect(e).To(gomega.BeNil())
			gomega.Expect(d.DrawID).To(gomega.Equal(1350))
		})

		ginkgo.It("should return empty/uninitialised draw", func() {
			d, e := ozLotto.GetDraw(10000)
			gomega.Expect(d.ID).To(gomega.Equal(0))
			gomega.Expect(d.Numbers).To(gomega.BeNil())
			gomega.Expect(e).NotTo(gomega.BeNil())
		})

		d, _ := ozLotto.GetDraw(1350)
		dt, _ := time.Parse("02/01/2006", "31/12/2019")

		ginkgo.It("should have these values", func() {
			gomega.Expect(d).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"ID":     gomega.Equal(100),
				"Game":   gomega.Not(gomega.BeNil()),
				"DrawID": gomega.Equal(1350),
				"Date":   gomega.Equal(dt),
			}))
		})

		ginkgo.It("should have these draw metrics", func() {
			gomega.Expect(d.Metrics.CrossSum).To(gomega.Equal(float64(118)))
			gomega.Expect(d.Metrics.CrossSumSupp).To(gomega.Equal(float64(200)))
			gomega.Expect(d.Metrics.Avg).To(gomega.Equal(16.86))
			gomega.Expect(d.Metrics.AvgSupp).To(gomega.Equal(22.22))

			//Using struct doesn't work. Don't know why yet.
			// gomega.Expect(&d.Metrics).To(gstruct.PointTo(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
			// 	"CrossSum":     gomega.Equal(float64(118)),
			// 	"CrossSumSupp": gomega.Equal(float64(200)),
			// 	"Avg":          gomega.Equal(16.86),
			// 	"AvgSupp":      gomega.Equal(22.22),
			// })))
		})
	})
})
