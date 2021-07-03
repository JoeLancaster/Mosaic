package sourceimage

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

type ColorImage SourceImage
type GrayImage SourceImage
type AverageType uint8

const (
	AVG_NONE   AverageType = 3
	AVG_MEAN   AverageType = 5
	AVG_MEDIAN AverageType = 7
	//default bins 16
	AVG_MODE    AverageType = 16
	AVG_MODE2   AverageType = 2
	AVG_MODE4   AverageType = 4
	AVG_MODE8   AverageType = 8
	AVG_MODE16  AverageType = 16
	AVG_MODE32  AverageType = 32
	AVG_MODE64  AverageType = 64
	AVG_MODE128 AverageType = 128
)

func div(num uint64, div uint64) uint64 {
	tmp := float64(num) / float64(div)
	return uint64(math.Floor(tmp))
}

type SourceImage struct {
	draw.Image
	Id      int
	Gray    bool
	average color.Color
}

type AverageableImg interface {
	draw.Image
	CalculateMean()
	Average() color.Color
}

func (s *SourceImage) Average() color.Color {
	return s.average
}

func (s *SourceImage) CalculateMean() {
	if s.Gray {
		tmp := GrayImage((*s))
		tmp.CalculateMean()
	} else {
		tmp := ColorImage((*s))
		tmp.CalculateMean()
	}
}

func (img *GrayImage) CalculateMean() {
	bounds := img.Bounds()
	var runningTot uint64 = 0
	Y := bounds.Dy()
	X := bounds.Dx()
	for y := 0; y < Y; y++ {
		for x := 0; x < X; x++ {
			tmp := draw.Image(img)
			runningTot += uint64(tmp.(*image.Gray).GrayAt(x, y).Y)

		}
	}
	img.average = color.Gray{uint8(div(runningTot, uint64(X*Y)))}

}

func (img *ColorImage) CalculateMean() {
	bounds := img.Bounds()
	runningTot := bigRGBA{0, 0, 0, 0}
	Y := bounds.Dy()
	X := bounds.Dx()
	count := uint64(X * Y)
	for y := 0; y < Y; y++ {
		for x := 0; x < X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			runningTot = bigRGBA{
				uint64(r) + runningTot.r,
				uint64(g) + runningTot.g,
				uint64(b) + runningTot.b,
				uint64(a) + runningTot.a}

		}
	}
	img.average = color.RGBA{
		uint8(div(runningTot.r, count)),
		uint8(div(runningTot.g, count)),
		uint8(div(runningTot.b, count)),
		uint8(div(runningTot.a, count)),
	}

}
