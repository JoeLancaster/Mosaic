package gray

import (
	"image"
	"image/color"
)

func Decolor(m *image.Image) *image.Gray {
	bounds := (*m).Bounds()
	result := image.NewGray(bounds)
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			px := (*m).At(x, y)
			np := color.GrayModel.Convert(px)
			result.Set(x, y, np)
		}
	}
	return result
}

func ConvertAll(imgs []*image.Image) []*image.Gray {
	var grays []*image.Gray
	for _, m := range imgs {
		g, ok := (*m).(*image.Gray)
		if !ok {
			g = Decolor(m)
		}
		grays = append(grays, g)
	}
	return grays
}
