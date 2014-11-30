package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/DAddYE/vips"
)

var (
	inputDir    = "/Users/cpjolicoeur/Dropbox/Wallpapers"
	outputDir   = "/tmp/image_resizing"
	vipsOptions = vips.Options{
		Width:        800,
		Height:       600,
		Crop:         false,
		Extend:       vips.EXTEND_WHITE,
		Interpolator: vips.BILINEAR,
		Gravity:      vips.CENTRE,
		Quality:      95,
	}
)

func scanDir(path string) (files []string, hello error) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, r := range entries {
		n := strings.ToUpper(r.Name())
		if strings.HasSuffix(n, ".JPG") || strings.HasSuffix(n, ".JPEG") || strings.HasSuffix(n, ".PNG") {
			files = append(files, path+"/"+r.Name())
		}
	}
	return
}

func batchResize(files []string, wg *sync.WaitGroup) {
	for i, origPath := range files {
		newPath := fmt.Sprintf("%s/thumb.%d.jpg", outputDir, i)
		go resizeImage(origPath, newPath, i, wg)
	}
}

func resizeImage(origName, newName string, num int, wg *sync.WaitGroup) {
	defer wg.Done()

	origFileStat, _ := os.Stat(origName)
	origFile, err := os.Open(origName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer origFile.Close()

	buf, err := ioutil.ReadAll(origFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	buf, err = vips.Resize(buf, vipsOptions)
	if err != nil {
		fmt.Println(err)
		return
	}

	cacheFile, err := os.Create(newName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cacheFile.Close()

	_, err = cacheFile.Write(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Resized image %d from %d to %d\n", num, origFileStat.Size(), int(len(buf)))
}

func main() {
	fmt.Println("Image resizer demo")

	runtime.GOMAXPROCS(runtime.NumCPU())

	files, _ := scanDir(inputDir)
	if len(files) == 0 {
		fmt.Println("no images found in", inputDir)
		return
	} else {
		fmt.Printf("found %d files in %s\n", len(files), inputDir)
	}

	var wg sync.WaitGroup
	wg.Add(len(files))

	batchResize(files, &wg)

	fmt.Println("Waiting to Finish")
	wg.Wait()
	fmt.Println("Terminating...")
}
