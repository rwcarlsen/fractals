package main

import "log"
import "gonum.org/v1/plot/plotter"
import "gonum.org/v1/plot"

func main() {
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}

	xys := plotter.XYs{{1, 2}, {3, 4}, {5, 6}}

	plotter.DefaultGliphStyle.Radius = vs.Points(2)

	scat, err := plotter.NewScatter(xys)
	if err != nil {
		log.Fatal(err)
	}
	p.Add(scatter)

	err := p.Save(width, height, "leaf.png")
	if err != nil {
		log.Fatal(err)
	}

}
