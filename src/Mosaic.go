package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"mosaic/stats/gray"
	"os"
	"path/filepath"
	"runtime"
)

func AvgWorker(jobs <-chan string, results chan<- color.RGBA, mins chan<- image.Rectangle, sourcePath string) {
	minx, miny := math.MaxInt32, math.MaxInt32
	for fname := range jobs {
		reader, err := os.Open(filepath.Join(sourcePath, fname))
		if err != nil {
			if fname != "" {
				fmt.Printf("Couldn't open: %s\n", fname)
			}
			fmt.Println(err)
			return
		}
		m, _, err := image.Decode(reader)
		if err != nil {
			fmt.Println(err)
			return
		}
		x, y := AbsDimension(m.Bounds())
		if x < minx {
			minx = x
		}
		if y < miny {
			miny = y
		}
		reader.Close()
		r, g, b, a := Average(m)
		//results <- color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
		//fmt.Printf("r: %d g: %d b: %d\n", r, g, b)
		results <- color.RGBA{r, g, b, a}
	}
	var b image.Rectangle
	b.Min.X = 0
	b.Min.Y = 0
	b.Max.X = minx
	b.Max.Y = miny
	mins <- b
	return
}

func main() {
	var targetFileName string
	var sourcePath string
	argLen := len(os.Args)
	if argLen < 3 { //too few arguments
		//print help message
	} else {
		targetFileName, _ = filepath.Abs(os.Args[1])
		sourcePath, _ = filepath.Abs(os.Args[2])
	}
	reader, err := os.Open(targetFileName)
	if err != nil {
		fmt.Println("Couldn't open target file.")
		fmt.Println(err)
		return
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	r, g, b, a := Average(m)
	fmt.Printf("Color averages in target image: r: %d, g: %d, b: %d, a: %d\n", r, g, b, a)
	files, err := ioutil.ReadDir(sourcePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	nofiles := len(files)
	results := make(chan color.RGBA, nofiles)
	imgs := make(chan string, nofiles)
	mins := make(chan image.Rectangle, runtime.NumCPU())

	fmt.Printf("Spawning: %d worker threads\n", runtime.NumCPU())
	for w := 1; w <= runtime.NumCPU(); w++ {
		go AvgWorker(imgs, results, mins, sourcePath)
	}
	for _, f := range files {
		imgs <- f.Name()
	}
	close(imgs)
	averages := make([]color.RGBA, nofiles)
	for r := 0; r < nofiles; r++ {
		col := <-results
		averages[r] = col //averages = append(averages, col)
		//		fmt.Printf("r: %d g: %d b: %d\n", col.R, col.G, col.B)
	}
	close(results)
	close(mins)
	var minD image.Rectangle
	minD.Min.X = 0
	minD.Min.Y = 0
	minD.Max.X = math.MaxInt32
	minD.Max.Y = math.MaxInt32
	for e := range mins {
		if e.Max.X < minD.Max.X {
			minD.Max.X = e.Max.X
		}
		if e.Max.Y < minD.Max.Y {
			minD.Max.Y = e.Max.Y
		}
	}
	fmt.Printf("min x: %d, min y: %d\n", minD.Max.X, minD.Max.Y)
	fmt.Printf("Finished averaging %d images.\n", nofiles)
	//bounds := m.Bounds()
	var new_bounds image.Rectangle
	dim := int(math.Floor(math.Sqrt(float64(nofiles))))
	new_bounds.Max.X = dim
	new_bounds.Max.Y = dim
	new_bounds.Min.X = 0
	new_bounds.Min.Y = 0
	new_img := image.NewRGBA(new_bounds)
	for x := new_bounds.Min.X; x < new_bounds.Max.X; x++ {
		for y := new_bounds.Min.Y; y < new_bounds.Max.Y; y++ {
			index := x + (dim * y)
			if index > nofiles {
				fmt.Printf("x: %d, y: %d\n", x, y)
				return
			}
			new_img.Set(x, y, averages[x+(new_bounds.Max.X*y)])
		}
	}
	new_img_file, err := os.Create("RESULT.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := png.Encode(new_img_file, new_img); err != nil {
		new_img_file.Close()
		fmt.Println(err)
		return
	}
	if err := new_img_file.Close(); err != nil {
		fmt.Println(err)
		return
	}
}
