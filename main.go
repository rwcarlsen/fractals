package main

import (
	"flag"
	"image/color"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func makeTransform(a, b, c, d, e, f float64) func(float64, float64) (float64, float64) {
	return func(x, y float64) (xout, yout float64) {
		xout = a*x + b*y + e
		yout = c*x + d*y + f
		return xout, yout
	}
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

type FastGlyph struct{}

func (FastGlyph) DrawGlyph(c *draw.Canvas, sty draw.GlyphStyle, pt vg.Point) {
	c.SetLineStyle(draw.LineStyle{Color: sty.Color, Width: vg.Points(float64(sty.Radius))})
	var p vg.Path
	p.Move(vg.Point{X: pt.X, Y: pt.Y})
	p.Line(vg.Point{X: pt.X + sty.Radius, Y: pt.Y + sty.Radius})
	c.Stroke(p)
}

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	pointsize := 0.3
	var height, width vg.Length = 700, 700
	n := 1000000

	pstem, stem := .01, makeTransform(0, 0, 0, .16, 0, 0)
	pleaflets, leaflets := .85, makeTransform(.85, .04, -.04, .85, 0, 1.6)
	pleftside, leftside := .07, makeTransform(.2, -.26, .23, .22, 0, 1.6)
	prightside, rightside := .07, makeTransform(-.15, .28, .26, .24, 0, .44)

	_ = pstem
	xys := make(plotter.XYs, n)
	x, y := 0.0, 0.0
	for i := 0; i < n; i++ {
		r := rand.Float64()
		fn := stem
		if r < pleaflets {
			fn = leaflets
		} else if r < pleaflets+pleftside {
			fn = leftside
		} else if r < pleaflets+pleftside+prightside {
			fn = rightside
		}
		x, y = fn(x, y)
		xys[i].X = x
		xys[i].Y = y
	}

	plotter.DefaultGlyphStyle.Radius = vg.Points(pointsize)
	plotter.DefaultGlyphStyle.Color = color.RGBA{0, 166, 0, 255}
	//plotter.DefaultGlyphStyle.Shape = draw.CircleGlyph{}
	//plotter.DefaultGlyphStyle.Shape = draw.BoxGlyph{}
	plotter.DefaultGlyphStyle.Shape = FastGlyph{}

	scat, err := plotter.NewScatter(xys)
	if err != nil {
		log.Fatal(err)
	}
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	p.Add(scat)

	p.HideAxes()
	err = p.Save(width, height, "leaf.svg")
	if err != nil {
		log.Fatal(err)
	}
	err = p.Save(width, height, "leaf.png")
	if err != nil {
		log.Fatal(err)
	}

}
