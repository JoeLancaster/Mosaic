package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"os"
	"runtime"
)

func AvgWorker(jobs <-chan string, results chan<- color.RGBA, source_file_path string) {
	for fname := range jobs {
		reader, err := os.Open(source_file_path + fname)
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
		reader.Close()
		r, g, b, a := Average(m)
		//results <- color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
		//fmt.Printf("r: %d g: %d b: %d\n", r, g, b)
		results <- color.RGBA{r, g, b, a}
	}
}

func main() {
	var target_file_name string
	var source_file_path string
	const NO_DELIM = "\nInput didn't end in a delimiter. Did you use C-d instead of RET?"
	argLen := len(os.Args)
	if argLen == 1 { //no arguments given
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter target image filename/path: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(NO_DELIM)
			return
		}
		target_file_name = text
		fmt.Print("Enter source images path: ")
		text, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println(NO_DELIM)
			return
		}
		source_file_path = text
	} else if argLen == 2 {
		fmt.Println("Run mosaic with two arguments. For example mosaic target_image /path/to/sources")
		return
	} else {
		target_file_name = os.Args[1]
		source_file_path = os.Args[2]
	}
	reader, err := os.Open(target_file_name)
	if err != nil {
		fmt.Println("Couldn't open target file.")
		fmt.Println(err)
		return
	}
	defer reader.Close()
	m, _, err := image.Decode(reader)
	r, g, b, a := Average(m)
	fmt.Printf("Color averages in target image: r: %d, g: %d, b: %d, a: %d\n", r, g, b, a)
	files, err := ioutil.ReadDir(source_file_path)
	if err != nil {
		fmt.Println(err)
		return
	}
	nofiles := len(files)
	results := make(chan color.RGBA, nofiles)
	imgs := make(chan string, nofiles)
	fmt.Printf("Getting image sizes.\n")
	var min_x, min_y uint32 = math.MaxUint32, math.MaxUint32
	for i := 0; i < nofiles; i++ {
		fname := files[i].Name()
		reader, err := os.Open(source_file_path + fname)
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
		if min_x > x {
			min_x = x
		}
		if min_y > y {
			min_y = y
		}
		reader.Close()
	}
	fmt.Printf("min x: %d, min_y: %d\n", min_x, min_y)
	fmt.Printf("Spawning: %d worker threads\n", runtime.NumCPU())
	for w := 1; w <= runtime.NumCPU(); w++ {
		go AvgWorker(imgs, results, source_file_path)
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