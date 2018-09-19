package main

import (
	"image/color"
	"log"
	"math"
	"os"

	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"

	"hepPlot/analysis"
)

func main() {

	// Generate Data
	nData := 1000
	//data := analysis.GenGaussStep(0., 1., 2., nData)
	data := analysis.GenLinGauss(25.e-3, 0., 0, 1., nData)

	// analysis
	lookback := 100
	smoothStrn := 0.95
	thresh := analysis.CalcThresholdSMA(data, lookback, smoothStrn)
	bias := analysis.CalcBias(data, thresh)

	// Create plot
	//TODO: need auto scaling axis here
	p1 := hplot.New()
	p1.Title.Text = "Time Series"
	p1.Y.Min = -5.
	p1.Y.Max = 5.
	p1.Y.Label.Text = "value"
	p1.X.Tick.Label.Color = color.RGBA{A: 0}

	// Draw threshold
	polyPts := make(plotter.XYs, 0, 2*nData)
	polyPts = append(polyPts, thresh.High...)
	//add points for low in backwards
	for ii := range thresh.Low {
		polyPts = append(polyPts, thresh.Low[len(thresh.Low)-ii-1])
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
	p1.Add(hplot.NewGrid())

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
