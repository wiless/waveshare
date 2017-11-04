package main

import (
	"image"
	"image/color"
	"os"

	"github.com/lucasb-eyer/go-colorful"

	"github.com/llgcode/draw2d/draw2dimg"

	"github.com/golang/glog"

	"golang.org/x/image/bmp"

	"github.com/wiless/waveshare"
)

var mono = true

func main() {
	waveshare.InitHW()
	var epd waveshare.EPD
	epd.SetDefaults()
	epd.Init(true)

	epdimg := ImageGenerate()
	epd.ClearFrame(0xff)
	epd.SetFrame(epdimg, 0, 0)
	epd.DisplayFrame()

}

func ImageGenerate() (epdimg image.Gray) {
	// img := image.NewGray(image.Rect(0, 0, 200, 210))
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))

	gc := draw2dimg.NewGraphicContext(img)
	gc.SetFillColor(color.Black)
	red, _ := colorful.Hex("#ff0000")
	gc.SetStrokeColor(red)
	gc.MoveTo(0, 0)
	gc.LineTo(150, 105)
	gc.QuadCurveTo(100, 20, 50, 20)
	gc.Close()
	gc.FillStroke()
	draw2dimg.SaveToPngFile("hello.png", img)
	/// grayimage
	b := img.Bounds()

	gimg := image.NewGray(b)
	var cg color.Gray

	for r := 0; r < b.Max.Y; r++ {
		for c := 0; c < b.Max.X; c++ {
			oldPixel := img.At(r, c)

			// gscale, _, _, _ := color.GrayModel.Convert(oldPixel).RGBA()
			cg = color.GrayModel.Convert(oldPixel).(color.Gray)

			// convert to monochrome
			if mono {
				if cg.Y > 0 {
					cg.Y = 255
				} else {
					cg.Y = 0
				}

			}
			gimg.SetGray(r, c, cg)

		}
	}

	epdimg = waveshare.Mono2ByteImage(gimg)

	f, e := os.Create("output.bmp")
	glog.Errorln(e)
	bmp.Encode(f, gimg)
	f.Close()

	return epdimg

}
