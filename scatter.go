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
	threshHigh := make(plotter.XYs, nData)
	threshLow := make(plotter.XYs, nData)
	bias := make(plotter.XYs, nData)
	for ii := range data {
		data[ii].X = float64(ii)
		data[ii].Y = dist.Rand()
		if ii < 20 {
			continue
		}
		threshHigh[ii].X = float64(ii)
		threshLow[ii].X = float64(ii)
		threshHigh[ii].Y = dist.Mu + 3.*dist.Sigma
		threshLow[ii].Y = dist.Mu - 3.*dist.Sigma
		if ii > 800 {
			threshHigh[ii].Y += (float64(ii) - 800.) / 100
			threshLow[ii].Y += (float64(ii) - 800.) / 100
		}
		bias[ii].X = float64(ii)
		center := (threshHigh[ii].Y + threshLow[ii].Y) / 2.
		width := (threshHigh[ii].Y - threshLow[ii].Y) / 2.
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

	// Draw threshold
	polyPts := make(plotter.XYs, 0, 2*nData)
	polyPts = append(polyPts, threshHigh...)
	//add points for low in backwards
	for ii := range threshLow {
		polyPts = append(polyPts, threshLow[len(threshLow)-ii-1])
	}

	poly, err := plotter.NewPolygon(polyPts)
	poly.Color = color.RGBA{B: 200, A: 20}
	poly.LineStyle.Color = color.RGBA{A: 0}

	// Draw data
	line, _, err := plotter.NewLinePoints(data)
	if err != nil {
		log.Panic(err)
	}
	line.Color = color.RGBA{G: 55, B: 200, R: 50, A: 255}
	p1.Add(line)

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
	lineBias, _, err := plotter.NewLinePoints(bias)
	if err != nil {
		log.Panic(err)
	}
	lineBias.Color = color.RGBA{G: 50, B: 155, R: 155, A: 255}
	p2.Add(lineBias)

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
	poly.Plot(p1.DataCanvas(top), p1.Plot)

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
