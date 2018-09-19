package analysis

import (
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot/plotter"
)

func GenLinGauss(slope, intercept, mu, sigma float64, nData int) plotter.XYs {
	gauss := GenGauss(mu, sigma, nData)
	for ii := range gauss {
		gauss[ii].Y += (intercept + gauss[ii].X*slope)
	}
	return gauss
}

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

func GenMem(nData int) plotter.XYs {
	// Create a normal distribution.
	dist := distuv.Uniform{
		Min: 0,
		Max: 1,
		Src: rand.New(rand.NewSource(0)),
	}
	base := 500000000 * dist.Rand()
	data := make(plotter.XYs, nData)
	for ii := 0; ii < nData; ii++ {
		data[ii].X = float64(ii)
		if ii < 4*nData/10 {
			data[ii].Y = base
		} else if (ii % 100) == 0 {
			base = 500000000 * dist.Rand()
			data[ii].Y = base
		} else {
			data[ii].Y = base
		}
	}
	return data
}
