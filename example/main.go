package main

import (
	"image"
	"log"
	"os"

	"github.com/golang/glog"

	"golang.org/x/image/bmp"

	"github.com/wiless/waveshare"
)

func main() {
	img := image.NewGray(image.Rect(0, 0, 16, 16))
	res := waveshare.Image2Byte(img)

	f, e := os.Create("output.bmp")
	glog.Errorln(e)
	bmp.Encode(f, img)
	f.Close()

	log.Print(res)
}
