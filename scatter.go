package main

import (
	"image/color"
	"log"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

var shadeColor = []color.Color{
	color.RGBA{R: 255, G: 255, B: 255, A: 255}, // white for the bottom
	color.RGBA{B: 200, A: 50},                  // actual shade color
}

func main() {

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Create data and highway
	nData := 1000
	data := make(plotter.XYs, nData)
	hwyHigh := make(plotter.XYs, nData)
	hwyLow := make(plotter.XYs, nData)
	for i := range data {
		data[i].X = float64(i)
		data[i].Y = dist.Rand()
		hwyHigh[i].X = float64(i)
		hwyLow[i].X = float64(i)
		hwyHigh[i].Y = dist.Mu + 3.*dist.Sigma
		hwyLow[i].Y = dist.Mu - 3.*dist.Sigma
	}

	// Create plot
	p, err := plot.New()
	if err != nil {
		log.Panic(err)
	}
	p.Title.Text = "Time Series"
	p.Y.Min = -5.
	p.Y.Max = 5.
	p.Y.Label.Text = "value"
	p.X.Label.Text = "time"

	// Draw high
	lineHigh, pointsHigh, err := plotter.NewLinePoints(hwyHigh)
	if err != nil {
		log.Panic(err)
	}

	lineHigh.Color = color.RGBA{A: 0}
	pointsHigh.Color = color.RGBA{A: 0}
	lineHigh.ShadeColor = &shadeColor[1]
	p.Add(lineHigh, pointsHigh)

	// Draw low
	lineLow, pointsLow, err := plotter.NewLinePoints(hwyLow)
	if err != nil {
		log.Panic(err)
	}

	lineLow.Color = color.RGBA{A: 0}
	pointsLow.Color = color.RGBA{A: 0}
	lineLow.ShadeColor = &shadeColor[0]
	p.Add(lineLow, pointsLow)

	// Draw data
	line, points, err := plotter.NewLinePoints(data)
	if err != nil {
		log.Panic(err)
	}

	line.Color = color.RGBA{G: 155, B: 155, R: 50, A: 255}
	points.Color = color.RGBA{A: 0}

	p.Add(line, points)

	// Print
	err = p.Save(50*vg.Centimeter, 25*vg.Centimeter, "testdata/timeseries.png")
	if err != nil {
		log.Panic(err)
	}

}
