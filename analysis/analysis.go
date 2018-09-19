package analysis

import (
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot/plotter"
)

type Threshold struct {
	High plotter.XYs
	Low  plotter.XYs
}

func CalcThresholdSMA(data plotter.XYs, lookback int, smoothStrn float64) Threshold {

	out := Threshold{High: make(plotter.XYs, len(data)), Low: make(plotter.XYs, len(data))}

	smoothAve := 0.
	smoothStd := 0.
	//calc over time
	for ii := range data {
		if ii < 20 {
			continue
		}
		// moving average model
		dataFloats := make([]float64, 0, lookback)
		start := ii - lookback
		if start < 0 {
			start = 0
		}
		for jj := start; jj < ii; jj++ {
			dataFloats = append(dataFloats, data[jj].Y)
		}
		ave, std := stat.MeanStdDev(dataFloats, nil)
		if ii == 20 {
			smoothAve = ave
			smoothStd = std
		} else {
			smoothAve = (1.-smoothStrn)*ave + smoothStrn*smoothAve
			smoothStd = (1.-smoothStrn)*std + smoothStrn*smoothStd
		}
		out.High[ii].X = float64(ii)
		out.High[ii].Y = smoothAve + 3.*smoothStd
		out.Low[ii].X = float64(ii)
		out.Low[ii].Y = smoothAve - 3.*smoothStd
	}
	return out
}

func CalcBias(data plotter.XYs, thresh Threshold) plotter.XYs {
	if len(data) != len(thresh.High) || len(data) != len(thresh.Low) {
		//err
		return nil
	}
	out := make(plotter.XYs, len(data))
	for ii := range data {
		if ii < 20 {
			continue
		}
		width := (thresh.High[ii].Y - thresh.Low[ii].Y) / 2.
		center := (thresh.High[ii].Y + thresh.Low[ii].Y) / 2.
		out[ii].X = float64(ii)
		if width > 1e-6 {
			out[ii].Y = (data[ii].Y - center) / width
		} else {
			out[ii].Y = (data[ii].Y - center)
		}
	}
	return out
}
