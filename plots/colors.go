package plots

import "image/color"

type barColor string

const (
	blue   barColor = "blue"
	red    barColor = "red"
	orange barColor = "orange"
	green  barColor = "green"
	purple barColor = "purple"
	cyan   barColor = "cyan"
)

var colorList = []barColor{blue, red, orange, green, purple, cyan}

var colors = map[barColor]color.RGBA{
	blue:   {R: 40, G: 90, B: 190, A: 255},
	red:    {R: 190, G: 50, B: 45, A: 255},
	orange: {R: 235, G: 130, B: 25, A: 255},
	green:  {R: 56, G: 145, B: 65, A: 255},
	purple: {R: 120, G: 30, B: 150, A: 255},
	cyan:   {R: 0, G: 155, B: 180, A: 255},
}
