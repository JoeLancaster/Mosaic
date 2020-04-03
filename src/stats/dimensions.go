package stats

import (
	"image"
)

func AbsDim(b image.Rectangle) (int, int) {
	return (b.Max.X - b.Min.X), (b.Max.Y - b.Min.Y)
}
