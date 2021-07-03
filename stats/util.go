package stats

import (
	"errors"
	"fmt"
	"github.com/JoeLancaster/mosaic/sourceimage"
	"image"
	"image/color"
	//	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
)

type CGPair struct {
	Gray    *image.Gray
	Col     *image.RGBA
	Id      int
	Average int
}

func CastGray(fileName string, i image.Image) (*image.Gray, error) {
	g, e := i.(*image.Gray)

	if !e {
		return nil, errors.New("Bad cast: " + fileName + " is not a gray image")
	}
	return g, nil
}

func BatchCastToGray(imgs []image.Image) []*image.Gray {
	var grays []*image.Gray
	i := 0
	for _, image := range imgs {
		casted, err := CastGray("", image)
		if err != nil {
			grays[i] = casted
			i++
		}
	}
	return grays
}

func OpenImage(fileName string) (image.Image, error) {
	r, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	m, _, err := image.Decode(r)
	//	g := m.(GrayImage)
	g := m.(*image.Gray)
	fmt.Println(g.At(0, 0))
	fmt.Println(g.GrayAt(0, 0))

	g.SetGray(0, 0, color.Gray{128})
	fmt.Println(g.GrayAt(0, 0))
	fmt.Println(reflect.TypeOf(m))
	return m, err
}

func OpenGray(fileName string) (*image.Gray, error) {
	r, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	m, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	g, e := m.(*image.Gray)
	if !e {
		return nil, errors.New("Bad cast: " + fileName + " is not a gray image")
	}
	r.Close()
	return g, nil

}

func OpenCol(fileName string) (*image.RGBA, error) {
	r, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	m, _, err := image.Decode(r)
	g, e := m.(*image.RGBA)
	if !e || err != nil {
		fmt.Printf("Bad convert: %s. Type: %s\n", fileName, reflect.TypeOf(m).String())
		return nil, err
	}
	r.Close()
	return g, err
}

func LoadGrays(srcPath string) []*image.Gray {
	files, _ := ioutil.ReadDir(srcPath)
	numFiles := len(files)
	imgs := make([]*image.Gray, numFiles)
	for i := 0; i < numFiles; i++ {
		image, _ := OpenGray(filepath.Join(srcPath, files[i].Name()))
		imgs[i] = image
	}
	return imgs
}

type PathPair struct {
	col  string
	gray string
}

func OpenSourceImage(path string, averageMethod sourceimage.AverageType) (*sourceimage.SourceImage, error) {
	var srcImg sourceimage.SourceImage
	r, err := os.Open(path)
	defer r.Close()
	if err != nil {
		return nil, err
	}
	m, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}
	switch m.(type) {
	case *image.Gray:
		srcImg.Gray = true
	default:
		srcImg.Gray = false
	}
	tmp := m.(sourceimage.AverageableImg)
	switch averageMethod {
	case sourceimage.AVG_NONE:
		return &srcImg, nil
	case sourceimage.AVG_MEAN:
		tmp.CalculateMean()
	case sourceimage.AVG_MEDIAN:
		//		avgr.CalculateMean()
	case sourceimage.AVG_MODE:
		//		avgr.CalculateMean()
	}
	return &srcImg, nil
}

func LoadBoth(srcGray string, srcCol string) []CGPair {
	grayFiles, _ := ioutil.ReadDir(srcGray)
	colFiles, _ := ioutil.ReadDir(srcCol)
	sort.Slice(colFiles, func(l int, r int) bool { return colFiles[l].Name() < colFiles[r].Name() })
	outBuf := []PathPair{}

	for _, g := range grayFiles {
		m := sort.Search(len(colFiles), func(i int) bool { return colFiles[i].Name() == g.Name() })
		if m > 0 {
			var pp PathPair
			pp.col = filepath.Join(srcCol, g.Name())
			pp.gray = filepath.Join(srcGray, g.Name())
			outBuf = append(outBuf, pp)
			//			fmt.Println("Send")
		}
	}

	length := len(outBuf)
	pairs := make(chan PathPair, length)
	results := make(chan CGPair, length)

	for w := 1; w <= runtime.NumCPU(); w++ {
		go func(pps <-chan PathPair, results chan<- CGPair, id int) {
			for p := range pps {
				//				fmt.Printf("[%d] Start\n", id)
				var ret CGPair
				ret.Gray, _ = OpenGray(p.gray)
				ret.Col, _ = OpenCol(p.col)
				//				if ret.gray != nil && ret.col != nil {
				//				fmt.Printf("Done %s\n", p.col)
				results <- ret
				//				fmt.Printf("[%d] Wait\n", id)

			}
			fmt.Printf("DONE: %d\n", id)
			return
		}(pairs, results, w)
	}

	for i := 0; i < length; i++ {
		pairs <- outBuf[i]
	}
	close(pairs)
	fmt.Println("CLOSED")
	cgpairs := make([]CGPair, length)
	for i := 0; i < length; i++ {
		cgpairs[i] = <-results
	}

	fmt.Println("DONE READ")
	close(results)
	return cgpairs
}

type arrayGetter func(s int) int

/*
** Choose a value from a sorted list with noise
** e.g
** [0, 0, 1, 1, 1, 2, 3, 3, 4, 4, 4, 5, 6]
    0  1  2  3  4  5  6  7  8  9 10 11 12
** start 2, noise 0 => return 5
** start 3 noise 1 => return a random num between 5 and 8
*/
func FuzzyPick(start int, noise int, get arrayGetter, max int, rnd rand.Rand) int {
	right := start + 1
	left := start
	for left >= 0 && get(left) >= start-noise {
		left--
	}
	for right < max && get(right) <= start+noise {
		right++
	}
	return rnd.Intn(right-left) + left

}

func LoadImages(fileNames <-chan string, imgs chan<- *image.Image, mins chan<- image.Rectangle) {
	minx, miny := math.MaxInt32, math.MaxInt32
	for f := range fileNames {
		m, err := OpenImage(f)
		if err == nil {
			x, y := AbsDim(m.Bounds())
			if x < minx {
				minx = x
			}
			if y < miny {
				miny = y
			}

			imgs <- &m
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
