package main

import (
	"image"
	"image/color"
	"os"
"log"
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
	epd.Init(true)

//	epdimg := ImageGenerate()
	epdimg:=waveshare.LoadImage("kavishbw.jpg")
	log.Println("Image is ",epdimg)	
	epd.ClearFrame(0x00)
	epd.SetFrame(*epdimg, 0, 0)

//	epd.WriteBytePixel(56,64,0x00,0Xaa,0x00)


	epd.DisplayFrame()

	epd.Sleep(true)
}

func ImageGenerate() (epdimg image.Gray) {

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
			oldPixel := img.At(c, r)

			// gscale, _, _, _ := color.GrayModel.Convert(oldPixel).RGBA()
			cg = color.GrayModel.Convert(oldPixel).(color.Gray)

			// convert to monochrome
			if mono {
				if cg.Y > 0 {
					cg.Y = 0
				} else {
					cg.Y = 0
				}

			}
			gimg.SetGray(c, r, cg)

		}
	}

	epdimg = waveshare.Mono2ByteImage(gimg)

	f, e := os.Create("output.bmp")
	glog.Errorln(e)
	bmp.Encode(f, gimg)
	f.Close()

	return epdimg

}
