package main

import (
	"fmt"
	"path"

	"github.com/JoeLancaster/mosaic/sourceimage"
	"github.com/JoeLancaster/mosaic/stats"

	//	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const usage = "%s: target_file /path/to/images [output_file]\n"
const defaultOutput = "output.jpg"

func gstr(gray bool) string {
	if gray {
		return "gray"
	}
	return "colour"
}

func main() {
	var outputArg string
	var averageMethod sourceimage.AverageType
	averageMethod = sourceimage.AVG_MEAN
	argLength := len(os.Args)
	if argLength < 2 {
		fmt.Println("Too few arguments. See --help")
		return
	}
	if argLength > 3 {
		outputArg = os.Args[3]
	} else {
		outputArg = defaultOutput
	}
	switch os.Args[1] {
	case "help", "--help", "-h", "-help", "--h":
		fmt.Printf(usage, path.Base(os.Args[0]))
		return
	}
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	targetPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	sourcePath, err := filepath.Abs(os.Args[2])
	if err != nil {
		log.Fatalln(err)
	}
	outputPath, err := filepath.Abs(outputArg)
	if err != nil {
		log.Fatalln(err)
	}
	outputFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer outputFile.Close()

	files, err := ioutil.ReadDir(sourcePath)
	if err != nil {
		log.Fatalln(err)
	}
	noOfFiles := len(files)

	loadedImages := make(chan *sourceimage.SourceImage, noOfFiles)
	fileNamesChan := make(chan string, noOfFiles)
	targetImage, err := stats.OpenSourceImage(targetPath, averageMethod)
	if err != nil {
		log.Printf("Couldn't open: %s\n", targetPath)
		log.Fatalln(err)
	} else {
		fmt.Printf("Using %s image %s\n", gstr(targetImage.Gray), targetPath)
	}
	fmt.Printf("Spawning %d threads for %d files\n", cpus, noOfFiles)
	//Open source images in parallel
	for w := 1; w <= cpus; w++ {
		go func(fileNames <-chan string, results chan<- *sourceimage.SourceImage) {
			for f := range fileNames {
				img, err := stats.OpenSourceImage(f, averageMethod)
				if err != nil {
					log.Println(err)
				}
				results <- img
			}
		}(fileNamesChan, loadedImages)
	}
	for i := 0; i < noOfFiles; i++ {
		fileNamesChan <- filepath.Join(sourcePath, files[i].Name())
	}
	close(fileNamesChan)
	var images []*sourceimage.SourceImage
	okays := 0
	//retrieve, discard failed images
	for i := 0; i < noOfFiles; i++ {
		got := <-loadedImages
		if got != nil {
			got.Id = i
			images = append(images, got)
			okays++
		}
	}
	delta := noOfFiles - okays
	if delta != 0 {
		log.Printf("Failed to open %d source images", delta)
	}
	// solution

	return

}
