package stats

import (
	"image"
)

func AbsDim(b image.Rectangle) (int, int) {
	return (b.Max.X - b.Min.X), (b.Max.Y - b.Min.Y)
}

func AspectRatio(b image.Rectangle) float32 {
	width, height := AbsDim(b)
	return float32(width / height)
}
