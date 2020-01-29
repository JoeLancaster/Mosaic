package main

import (
	//	"fmt"
	"image"
)

type RGBA struct {
	r uint32
	g uint32
	b uint32
	a uint32
}

func Average(m image.Image) (uint8, uint8, uint8, uint8) {
	bounds := m.Bounds()
	var cr, cg, cb, ca uint32
	var rr, rg, rb, ra uint32 = 0, 0, 0, 0
	xdif, ydif := AbsDimension(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		cr, cg, cb, ca = 0, 0, 0, 0
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := m.At(x, y).RGBA()
			if y == 0 && x == 0 {

			}
			cr += (r >> 8)
			cg += (g >> 8)
			cb += (b >> 8)
			ca += (a >> 8)
		}
		rr += (cr / xdif)
		rg += (cg / xdif)
		rb += (cb / xdif)
		ra += (ca / xdif)
	}
	rr /= ydif
	rg /= ydif
	rb /= ydif
	ra /= ydif
	//	fmt.Printf("r: %d g: %d b: %d\n", rr, rg, rb)
	return uint8(rr), uint8(rg), uint8(rb), uint8(ra)
}

func AspectRatio(b image.Rectangle) float32 {
	width, height := AbsDimension(b)
	return float32(width / height)
}

func AbsDimension(b image.Rectangle) (uint32, uint32) {
	xdif := uint32(b.Max.X - b.Min.X)
	ydif := uint32(b.Max.Y - b.Min.Y)
	return xdif, ydif
}
