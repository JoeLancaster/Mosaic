package solutions

import (
	"image"
	"os"
	"runtime"
)

func OpenGray(fileName string) (*image.Gray, error) {
	r, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	m, _, err := image.Decode(r)
	g, e := m.(*image.Gray)
	if !e {
		errors.New("Bad cast: " + fileName + " is not a gray image")
	}
	r.Close()
	return g, err

}
