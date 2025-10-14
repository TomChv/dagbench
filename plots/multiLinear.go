package plots

import (
	"image/color"
	"math"
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// MultiLinearByVersion creates a plot with multiple lines
// for each operations with version as X axis.
//
//nolint:funlen
func MultiLinearByVersion(data DataSet, title string) (*plot.Plot, error) {
	p := plot.New()
	p.Title.Text = title
	p.Title.TextStyle.Font.Size = vg.Points(22)
	p.X.Label.Text = "Version"
	p.X.Label.TextStyle.Font.Size = font.Points(12)
	p.Y.Label.Text = defaultYLabel
	p.Y.Min = 0
	p.Y.Label.TextStyle.Font.Size = font.Points(12)
	p.Add(plotter.NewGrid())

	// Get the versions from the data
	var versions []string
	if len(versions) == 0 {
		// derive from data if not provided
		vset := map[string]struct{}{}
		for _, m := range data {
			for v := range m {
				vset[v] = struct{}{}
			}
		}
		for v := range vset {
			versions = append(versions, v)
		}
		sort.Strings(versions)
	}

	// Add X axis tags
	p.NominalX(versions...)
	p.X.Min = -0.1
	p.X.Max = float64(len(versions)-1) + 0.1
	p.X.Tick.Label.Rotation = math.Pi / 6 // 30°
	p.X.Tick.Label.Font.Size = font.Points(10)

	// Deterministic order of operations in legend
	opNames := make([]string, 0, len(data))
	for op := range data {
		opNames = append(opNames, op)
	}
	sort.Strings(opNames)

	// Add lines to plot
	for i, op := range opNames {
		pts := make(plotter.XYs, len(versions))
		for j, v := range versions {
			pts[j].X = float64(j) // index on nominal axis

			if val, ok := data[op][v]; ok {
				pts[j].Y = val
			} else {
				pts[j].Y = math.NaN() // break the line if missing
			}
		}
		line, dots, err := plotter.NewLinePoints(pts)
		if err != nil {
			return nil, err
		}
		c := colors[colorList[i%len(colorList)]]
		line.Color = c
		line.Width = vg.Points(2)
		dots.Radius = vg.Points(2.5)
		dots.Color = c

		p.Add(line, dots)
		p.Legend.Add(op, line, dots)
	}

	// Customize plot styles
	p.Legend.TextStyle.Font.Size = font.Points(10)
	p.Legend.Padding = vg.Points(6)
	p.Legend.ThumbnailWidth = vg.Points(14)
	p.Legend.TextStyle.Font.Size = vg.Points(11)
	p.Legend.TextStyle.Color = color.Black
	p.Legend.Top = true

	return p, nil
}
