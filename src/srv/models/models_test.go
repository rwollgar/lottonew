package models

import (
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

		record := []string{
			"1350", "31/12/2019", "10", "11", "29", "1", "39", "23", "5",
			"42", "40", "0", "0", "0.00", "0.00", "1", "0", "123427.60",
			"123427.60", "66", "11", "254116.50", "3850.25", "343", "35",
			"130683.00", "381.00", "3864", "461", "152434.80", "39.45",
			"79375", "9827", "1742281.25", "21.95", "118330", "14958", "1952445.00",
			"16.50"}

		g := game{
			Name:             "oz-lotto",
			StandardNumbers:  7,
			Supplementary:    2,
			MaxNumber:        45,
			MaxSupplementary: 45,
			Format:           "csv",
			DrawOffset:       1}

		d := initDraw(draw{Game: &g, ID: 1}, record)
		dt, _ := time.Parse("02/01/2006", "31/12/2019")

		ginkgo.It("should have these values", func() {
			gomega.Expect(d).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"ID":     gomega.Equal(1),
				"Game":   gomega.Not(gomega.BeNil()),
				"DrawID": gomega.Equal(1350),
				"Date":   gomega.Equal(dt),
			}))
		})

		// CrossSum:118
		// Avg:16.857142857142858
		// CrossSumSupp:200
		// AvgSupp:22.22222222222222
		ginkgo.It("should have these draw metrics", func() {
			gomega.Expect(d.Metrics).To(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"CrossSum":     gomega.Equal(118.0),
				"CrossSumSupp": gomega.Equal(200.0),
				"Avg":          gomega.Equal(16.857142857142858),
				"AvgSupp":      gomega.Equal(22.22222222222222),
			}))
		})
	})
})
