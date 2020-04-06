package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/JoeLancaster/mosaic/stats"
)

func main() {

	grayFlag := true

	argLen := len(os.Args)
	if argLen < 2 {
		fmt.Println("Too few arguments. Run mosaic with at least two arguments. Run mosaic --help for help")
		return
	}
	switch os.Args[1] {
	case "help", "--help", "-h", "-help", "--h":
		fmt.Println("Run mosaic with two arguments: the target image and the source directory of images.\n\tEx: mosaic house.jpg ~/Pictures/Trees")
		return
	}

	targetFileName, err := filepath.Abs(os.Args[1])
	sourcePath, err := filepath.Abs(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	m, err := stats.OpenImage(targetFileName)

	if err != nil {
		log.Fatal(err)
	}
	files, err := ioutil.ReadDir(sourcePath)
	if err != nil {
		log.Fatal(err)
	}
	nofiles := len(files)
	images := make(chan string, nofiles)

	numCores := runtime.NumCPU()

	var grayResults chan uint8
	var colResults chan color.RGBA

	if grayFlag {
		grayResults = make(chan uint8, nofiles)
	} else {
		colResults = make(chan color.RGBA, nofiles)
	}

	fmt.Printf("Spawning %d worker threads\n", numCores)

	for w := 1; w <= numCores; w++ {
		if grayFlag {
			//go gray worker
		} else {
			//go color worker
		}
	}

}
