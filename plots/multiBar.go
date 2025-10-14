package plots

import (
	"image/color"
	"math"
	"os"
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgsvg"
)

func SaveWithOutsideLegend(
	p *plot.Plot,
	totalWidth, height, legendWidth vg.Length,
	filename string,
) error {
	svg := vgsvg.New(totalWidth, height)
	dc := draw.New(svg)

	// Create a white background SVG canvas.
	var bg vg.Path
	bg.Move(vg.Point{X: 0, Y: 0})
	bg.Line(vg.Point{X: totalWidth, Y: 0})
	bg.Line(vg.Point{X: totalWidth, Y: height})
	bg.Line(vg.Point{X: 0, Y: height})
	bg.Close()
	dc.SetColor(color.White)
	dc.Fill(bg)

	// Disable legend for the original plot to draw it manually after
	originalLegend := p.Legend
	p.Legend = plot.NewLegend()

	// Draw the original plot area on the left
	plotWidth := totalWidth - legendWidth - vg.Points(20) // leave a 30 pt gap for the legend
	dcPlot := draw.Crop(dc, vg.Points(0), -(totalWidth - plotWidth), vg.Points(0), vg.Points(0))
	p.Draw(dcPlot)

	// Restore legend and style it.
	p.Legend = originalLegend
	p.Legend.Padding = vg.Points(6)
	p.Legend.ThumbnailWidth = vg.Points(14)
	p.Legend.TextStyle.Font.Size = vg.Points(11)
	p.Legend.TextStyle.Color = color.Black

	// Draw legend manually on the right
	dc.Push()
	dc.Translate(vg.Point{X: -1 * vg.Inch, Y: -0.5 * vg.Inch})
	p.Legend.Draw(dc)
	dc.Pop()

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close() //nolint:errcheck

	if _, err := svg.WriteTo(f); err != nil {
		return err
	}

	return nil
}

//nolint:funlen
func MultiBar(dataset DataSet, title string) (*plot.Plot, error) {
	p := plot.New()

	// Collect all unique operation names
	nameSet := make(map[string]struct{})
	for _, data := range dataset {
		for name := range data {
			nameSet[name] = struct{}{}
		}
	}
	names := make([]string, 0, len(nameSet))
	for name := range nameSet {
		names = append(names, name)
	}
	sort.Strings(names)

	p.NominalX(names...)

	groupWidth := vg.Points(50)
	barWidth := groupWidth / vg.Length(len(dataset)+1)
	offset := -groupWidth / 2

	i := 0
	for label, data := range dataset {
		values := make(plotter.Values, len(names))
		for j, name := range names {
			values[j] = data[name]
		}

		bc, err := plotter.NewBarChart(values, barWidth)
		if err != nil {
			return nil, err
		}

		bc.LineStyle.Width = 0
		bc.Color = colors[colorList[i%len(colorList)]]
		bc.Offset = offset + barWidth*vg.Length(i)
		p.Add(bc)
		p.Legend.Add(label, bc)
		i++
	}

	// Customize plot styles
	p.Legend.Top = true
	p.Legend.XOffs = vg.Points(60)
	p.Legend.YOffs = vg.Points(15)
	p.Legend.TextStyle.Font.Size = font.Points(11)
	p.Legend.Padding = vg.Points(4)
	p.Legend.ThumbnailWidth = vg.Points(12)

	p.Title.Text = title
	p.Title.TextStyle.Font.Size = vg.Points(22)
	p.X.Label.Text = defaultXLabel
	p.X.Label.TextStyle.Font.Size = font.Points(12)
	p.Y.Label.Text = defaultYLabel
	p.Y.Label.TextStyle.Font.Size = font.Points(12)
	p.X.Padding = vg.Points(15)
	p.Y.Padding = vg.Points(5)

	// Rotate x labels to fit better
	p.X.Tick.Label.Rotation = math.Pi / 4 // 45 degrees
	p.X.Tick.Label.XAlign = text.XRight
	p.X.Tick.Label.YAlign = text.YCenter
	p.X.Tick.Label.Font.Size = font.Points(11)

	// Add grid for readability
	p.Add(plotter.NewGrid())

	return p, nil
}
