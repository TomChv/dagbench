package plots

import (
	"math"
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"
)

// SingleBar creates a plot with a single bar for each operations.
func SingleBar(data Data, title string) (*plot.Plot, error) {
	p := plot.New()

	// Sort names for consistent ordering
	names := make([]string, 0, len(data))
	for name := range data {
		names = append(names, name)
	}
	sort.Strings(names)

	// Little hack to center the bars
	names = append([]string{""}, append(names, "")...)

	values := make(plotter.Values, len(names))
	for i, name := range names {
		if name == "" {
			values[i] = 0

			continue
		}

		values[i] = data[name]
	}

	barWidth := vg.Points(14)
	bars, err := plotter.NewBarChart(values, barWidth)
	if err != nil {
		return nil, err
	}

	// Customize bar styles
	bars.LineStyle.Width = vg.Length(0)
	bars.Color = colors[blue]
	bars.Width = barWidth
	bars.Offset = vg.Points(0)

	p.Add(bars)
	p.NominalX(names...)

	// Customize plot styles
	p.Title.Text = title
	p.Title.TextStyle.Font.Size = vg.Points(22)
	p.X.Label.Text = defaultXLabel
	p.X.Label.TextStyle.Font.Size = font.Points(12)
	p.Y.Label.Text = defaultYLabel
	p.Y.Label.TextStyle.Font.Size = font.Points(12)
	p.X.Padding = vg.Points(15)

	// Rotate x labels to fit better
	p.X.Tick.Label.Rotation = math.Pi / 4 // 45 degrees
	p.X.Tick.Label.XAlign = text.XRight
	p.X.Tick.Label.YAlign = text.PosCenter
	p.X.Tick.Label.Font.Size = font.Points(11)

	// Add grid for readability
	p.Add(plotter.NewGrid())

	return p, nil
}
