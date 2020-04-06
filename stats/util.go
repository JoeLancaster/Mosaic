package stats

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
)

func IsPowerOfTwo(n uint8) bool {
	return (n & (n - 1)) == 0
}

func OpenImage(fileName string) (*image.Image, error) {
	r, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	m, _, err := image.Decode(r)

	if err != nil {
		return nil, err
	}
	r.Close()
	return &m, nil
}

func LoadImages(fileNames <-chan string, imgs chan<- *image.Image, mins chan<- image.Rectangle) {
	minx, miny := math.MaxInt32, math.MaxInt32
	for f := range fileNames {
		m, err := OpenImage(f)
		if err == nil {
			x, y := AbsDim((*m).Bounds())
			if x < minx {
				minx = x
			}
			if y < miny {
				miny = y
			}

			imgs <- m
		} else {
			log.Printf("File %s not used.\nErr:\t%s\n", f, err)
		}
	}
	var b image.Rectangle
	b.Min.X = 0
	b.Min.Y = 0
	b.Max.X = minx
	b.Max.Y = miny
	mins <- b
	return
}
