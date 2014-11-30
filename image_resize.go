package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

func batchResize(files []string) {
	for i, origPath := range files {
    fmt.Printf("File %d: ", i+1)
		newPath := fmt.Sprintf("%s/thumb.%d.jpg", outputDir, i)
		newSize, oldSize := resizeImage(origPath, newPath)
    fmt.Printf("Resized from %d to %d\n", oldSize, newSize)
	}
}

func resizeImage(origName, newName string) (int, int64) {
	// fmt.Println("vips options:", vipsOptions)
	origFileStat, _ := os.Stat(origName)
	origFile, err := os.Open(origName)
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}
	defer origFile.Close()

  buf, err := ioutil.ReadAll(origFile)
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}

	buf, err = vips.Resize(buf, vipsOptions)
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}

  cacheFile, err := os.Create(newName)
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}
  defer cacheFile.Close()

  _, err = cacheFile.Write(buf)
	if err != nil {
		fmt.Println(err)
		return 0, origFileStat.Size()
	}

	return int(len(buf)), origFileStat.Size()
}

func main() {
	fmt.Println("Image resizer demo")

	files, _ := scanDir(inputDir)
	if len(files) == 0 {
		fmt.Println("no images found in", inputDir)
		return
	}

  batchResize(files)
	// for _, file := range files {
	// 	resizeImage(file)
	// }
	// fmt.Println("image files found:", files)
}
