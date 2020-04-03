package gray

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "math"
	//	"math/bits"
	//	"math/rand"
	"os"
	"testing"
	//	"time"
)

func mkSolid(L uint8) *image.Gray {
	var bounds image.Rectangle
	bounds.Min.X = 0
	bounds.Min.Y = 0
	bounds.Max.X = 64
	bounds.Max.Y = 64
	testImg := image.NewGray(bounds)
	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			testImg.SetGray(x, y, color.Gray{Y: L})
		}
	}
	return testImg
}

func TestMean(t *testing.T) {
	peppersFile, err := os.Open("peppers-gray.jpg")
	defer peppersFile.Close()

	if err != nil {
		t.Errorf("Problem opening test image \"peppers-gray.jpg")
	}

	mi, _, err := image.Decode(peppersFile)
	m := mi.(*image.Gray)

	ans := mean(m)
	fmt.Printf("peppers mean luminosity: %d\n", ans)
	if ans > 255 || ans < 0 {
		t.Errorf("Mean(m) got: %d", ans)
	}

	blockWhite := mkSolid(255)
	blockBlack := mkSolid(0)

	ansW := Mean(blockWhite)
	ansB := Mean(blockBlack)
	if ansW != 255 {
		t.Errorf("Mean of a pure white image got: %d", ansW)
	}
	if ansB != 0 {
		t.Errorf("Mean of a pure black image got: %d", ansB)
	}

}

// func TestExtractMode(t *testing.T) {
// 	rand.Seed(time.Now().UnixNano())
// 	for bins := 1; bins < 256; bins <<= 1 {
// 		shift := uint8(bits.LeadingZeros8(uint8(bins))) + 1
// 		hist := make([]int, bins)
// 		for i, _ := range hist {
// 			hist[i] = rand.Int()
// 		}
// 		ans := extractMode(hist, shift)
// 		fmt.Printf("Random mode with %d bins: %d\n", bins, ans)
// 	}
// }

func TestBinnedMode(t *testing.T) {
	peppersFile, err := os.Open("peppers-gray.jpg")
	defer peppersFile.Close()

	if err != nil {
		t.Errorf("Problem opening test image \"peppers-gray.jpg")
	}

	mi, _, err := image.Decode(peppersFile)
	m := mi.(*image.Gray)

	for bins := 1; bins < 256; bins <<= 1 {
		bin8 := uint8(bins)
		ans := BinnedMode(m, bin8)
		fmt.Printf("Peppers mode with %d bins: %d\n", bins, ans)
	}

	testImg := mkSolid(255)

	i := 128
	expV := 0
	for bins := 1; bins < 256; bins <<= 1 {
		expV = expV | i
		bin8 := uint8(bins)
		ans := BinnedMode(testImg, bin8)
		if ans != uint8(expV) {
			t.Errorf("Binned mode pure white expected: %d, got %d", expV, ans)
		}
		fmt.Printf("Pure white mode with %d bins: %d\n", bins, ans)
		i >>= 1
	}

}
