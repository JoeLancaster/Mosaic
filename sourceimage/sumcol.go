package sourceimage

import (
	"image/color"
	"math"
)

func div32(num uint64, div uint32) uint64 {
	tmp := float64(num) / float64(div)
	return uint64(math.Floor(tmp))
}

type bigRGBA struct {
	r uint64
	g uint64
	b uint64
	a uint64
}

type bigGray uint64

type SummableColor interface {
	Add(s color.Color) SummableColor
	Divide(count uint32) SummableColor
}

func (x bigGray) Divide(count uint32) SummableColor {
	return bigGray(div32(uint64(x), count))
}

func (x bigGray) Add(tmp color.Color) SummableColor {
	b := tmp.(color.Gray)
	return bigGray(uint64(x) + uint64(b.Y))
}

func (x bigRGBA) Add(c color.Color) SummableColor {
	r, g, b, a := c.RGBA()
	return bigRGBA{
		x.r + uint64(r),
		x.g + uint64(g),
		x.b + uint64(b),
		x.a + uint64(a)}
}

func (x bigRGBA) Divide(count uint32) SummableColor {
	return bigRGBA{
		div32(x.r, count),
		div32(x.g, count),
		div32(x.b, count),
		div32(x.a, count)}
}
