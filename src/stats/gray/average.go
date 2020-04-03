package gray

import (
	"image"
	"math"
	"math/bits"
	"mosaic/stats"
)

func extractMode(hist []int, shift uint8) uint8 {
	msf := 0
	msfAt := 0
	for i, v := range hist {
		if v >= msf { //>= if contiguous bins are the same value then prefer the brighter bin
			msf = v
			msfAt = i
		}
	}
	//binned modal is between i << shift and (i + 1) << shift - 1
	//so take the difference and use that as the mode
	tmp := (msfAt << shift) + ((msfAt + 1) << shift)
	return uint8(tmp >> 1) //fast divide by two
}

func BinnedMode(m *image.Gray, bins uint8) uint8 {
	bounds := m.Bounds()
	X, Y := stats.AbsDim(bounds)

	hist := make([]int, bins)

	shift := uint8(bits.LeadingZeros8(bins)) + 1

	for x := 0; x < X; x++ {
		for y := 0; y < Y; y++ {
			pv := m.GrayAt(x, y).Y
			hist[pv>>shift]++
		}
	}
	return extractMode(hist, shift)
}

func mean(m *image.Gray) int {
	runningTot := 0
	bounds := m.Bounds()
	X, Y := stats.AbsDim(bounds)
	numPix := X * Y

	for x := 0; x < X; x++ {
		for y := 0; y < Y; y++ {
			runningTot += int(m.GrayAt(x, y).Y)
		}
	}

	return int(math.Floor(float64(runningTot) / float64(numPix)))
}

func Mean(m *image.Gray) uint8 {
	return uint8(mean(m))
}
