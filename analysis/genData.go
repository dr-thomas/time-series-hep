package analysis

import (
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot/plotter"
)

func GenGauss(mu, sigma float64, nData int) plotter.XYs {
	return GenGaussStep(mu, sigma, 0., nData)
}

func GenGaussStep(mu, sigma, step float64, nData int) plotter.XYs {
	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    mu,
		Sigma: sigma,
		Src:   rand.New(rand.NewSource(0)),
	}
	// Create data and highway
	data := make(plotter.XYs, nData)
	for ii := range data {
		data[ii].X = float64(ii)
		data[ii].Y = dist.Rand()
		if ii > len(data)*3/4 {
			if mu > 1e-6 {
				data[ii].Y += mu * step
			} else {
				data[ii].Y += step
			}
		}
	}
	return data
}
