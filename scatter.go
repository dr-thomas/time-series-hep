package main

import (
	"image/color"
	"log"
	"math"
	"os"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"

	"go-hep.org/x/hep/hplot"
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
	bias := make(plotter.XYs, nData)
	for ii := range data {
		data[ii].X = float64(ii)
		data[ii].Y = dist.Rand()
		hwyHigh[ii].X = float64(ii)
		hwyLow[ii].X = float64(ii)
		hwyHigh[ii].Y = dist.Mu + 3.*dist.Sigma
		hwyLow[ii].Y = dist.Mu - 3.*dist.Sigma
		if ii > 800 {
			hwyHigh[ii].Y += (float64(ii) - 800.) / 100
			hwyLow[ii].Y += (float64(ii) - 800.) / 100
		}
		bias[ii].X = float64(ii)
		center := (hwyHigh[ii].Y + hwyLow[ii].Y) / 2.
		width := (hwyHigh[ii].Y - hwyLow[ii].Y) / 2.
		if width > 0. {
			bias[ii].Y = (data[ii].Y - center) / width
		} else {
			bias[ii].Y = (data[ii].Y - center)
		}
	}

	// Create plot
	p1 := hplot.New()
	p1.Title.Text = "Time Series"
	p1.Y.Min = -5.
	p1.Y.Max = 5.
	p1.Y.Label.Text = "value"
	p1.X.Tick.Marker = hplot.NoTicks{}

	// Draw high
	lineHigh, pointsHigh, err := plotter.NewLinePoints(hwyHigh)
	if err != nil {
		log.Panic(err)
	}

	lineHigh.Color = color.RGBA{A: 0}
	pointsHigh.Color = color.RGBA{A: 0}
	lineHigh.ShadeColor = &shadeColor[1]
	p1.Add(lineHigh, pointsHigh)

	// Draw low
	lineLow, pointsLow, err := plotter.NewLinePoints(hwyLow)
	if err != nil {
		log.Panic(err)
	}

	lineLow.Color = color.RGBA{A: 0}
	pointsLow.Color = color.RGBA{A: 0}
	lineLow.ShadeColor = &shadeColor[0]
	p1.Add(lineLow, pointsLow)

	// Draw data
	line, points, err := plotter.NewLinePoints(data)
	if err != nil {
		log.Panic(err)
	}

	line.Color = color.RGBA{G: 155, B: 155, R: 50, A: 255}
	points.Color = color.RGBA{A: 0}

	p1.Add(line, points)

	// Create lower plot
	p2 := hplot.New()
	p2.Y.Label.Text = "bias"
	p2.X.Label.Text = "time"
	p2.Add(hplot.NewGrid())

	//Determine bias mixima
	min := 100.
	max := -100.
	for _, xx := range bias {
		if xx.Y < min {
			min = xx.Y
		}
		if xx.Y > max {
			max = xx.Y
		}
	}
	p2.Y.Min = min
	p2.Y.Max = max

	// Draw bias
	lineBias, pointsBias, err := plotter.NewLinePoints(bias)
	if err != nil {
		log.Panic(err)
	}

	lineBias.Color = color.RGBA{G: 50, B: 155, R: 155, A: 255}
	pointsBias.Color = color.RGBA{A: 0}

	p2.Add(lineBias, pointsBias)

	const (
		width  = 40 * vg.Centimeter
		height = width / math.Phi
	)

	c := vgimg.PngCanvas{Canvas: vgimg.New(width, height)}
	dc := draw.New(c)
	top := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0, Y: 0.3 * height},
			Max: vg.Point{X: width, Y: height},
		},
	}
	p1.Draw(top)

	bottom := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0, Y: 0},
			Max: vg.Point{X: width, Y: 0.3 * height},
		},
	}
	p2.Draw(bottom)

	f, err := os.Create("testdata/timeseries_bias.png")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	defer f.Close()
	_, err = c.WriteTo(f)
	if err != nil {
		log.Fatal(err)
	}
	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

}
