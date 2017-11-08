package main

import (
	"image"
	"image/color"
	"os"
	"time"

	"github.com/llgcode/draw2d"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"

	"github.com/golang/freetype/truetype"
	"github.com/golang/glog"

	"golang.org/x/image/bmp"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/wiless/waveshare"
)

var mono = true
var epd waveshare.EPD

func main() {
	waveshare.InitHW()
	draw2d.SetFontFolder(".")

	epdimg := ImageGenerate()
	UpdateImage(epdimg)
	time.Sleep(500 * time.Millisecond)
	PartialUpdate()
}
func UpdateImage(epdimg image.Gray) {
	epd.Init(true)
	epd.ClearFrame(0xff)
	epd.SetFrame(epdimg)
	epd.DisplayFrame()

	// epd.ClearFrame(0xff)
	// epd.SetFrame(epdimg)
	// epd.DisplayFrame()
}

func PartialUpdate() {

	epd.Init(false)
	timeimg := image.NewRGBA(image.Rect(0, 0, 100, 48))
	gc := draw2dimg.NewGraphicContext(timeimg)
	gc.ClearRect(0, 0, 100, 30)
	gc.SetFillColor(color.Black)
	gc.SetLineWidth(1.5)
	gc.StrokeStringAt("Hey I am good", 0, 10)
	gc.FillStroke()
	gc.Save()
	draw2dimg.SaveToPngFile("subimage.png", timeimg)
	gimg := ConvertToGray(timeimg)
	SaveBMP("subimage.bmp", gimg)
	epd.SetSubFrame(8, 8, gimg)

}
func ConvertToGray(cimg image.Image) *image.Gray {
	b := cimg.Bounds()
	gimg := image.NewGray(b)
	var cg color.Gray
	mono = true
	for r := 0; r < b.Max.Y; r++ {
		for c := 0; c < b.Max.X; c++ {
			oldPixel := cimg.At(r, c)

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
	return gimg
}

func SaveBMP(fname string, img image.Image) {
	fp, fe := os.Create(fname)
	if fe != nil {
		glog.Errorln("Unable to Save ", fname)
		return
	}
	bmp.Encode(fp, img)
	fp.Close()
}

func ImageGenerate() (epdimg image.Gray) {
	// img := image.NewGray(image.Rect(0, 0, 200, 210))
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	// for r := 0; r < 200; r++ {
	// 	for c := 0; c < 200; c++ {
	// 		img.Set(r, c, color.White)
	// 	}
	// }

	gc := draw2dimg.NewGraphicContext(img)
	gc.ClearRect(0, 0, 200, 200)

	gc.Save()
	gc.SetStrokeColor(color.Black)
	gc.SetFillColor(color.Black)
	gc.SetLineWidth(2)
	draw2dkit.Rectangle(gc, 30, 30, 100, 100)
	gc.Stroke()
	gc.Save()

	gc.SetStrokeColor(color.Black)
	gc.SetFillColor(color.White)
	gc.SetLineWidth(4)
	draw2dkit.Circle(gc, 100, 100, 30)
	gc.FillStroke()

	draw2dkit.RoundedRectangle(gc, 105, 105, 180, 180, 10, 10)
	gc.Stroke()

	gc.SetFillColor(color.Black)

	gc.SetStrokeColor(color.Black)
	// gc.Close()
	// gc.Restore()
	// gc.SetFillColor(color.Black)
	font, _ := truetype.Parse(goregular.TTF)
	// font, _ := truetype.Parse(gobold.TTF)

	gc.SetFont(font)
	gc.SetFontSize(14)
	gc.SetLineWidth(2.5)
	msg := "ABCDEFGHIJKLMNOPQ"
	// L, T, R, B := gc.GetStringBounds(msg)

	// fmt.Println("L T R B", L, T, R, B)
	gc.StrokeStringAt(msg, 0, 20)
	// gc.FillStroke()
	gc.SetFontSize(20)
	gc.SetLineWidth(4)
	datestr := time.Now().Format(time.Stamp)
	gc.StrokeStringAt(datestr, 0, 170)
	gc.FillStroke()
	gc.Save()

	gc.Save()
	gc.Close()
	draw2dimg.SaveToPngFile("hello.png", img)

	f1, _ := os.Create("input.bmp")

	bmp.Encode(f1, img)
	f1.Close()

	/// grayimage
	b := img.Bounds()

	gimg := image.NewGray(b)
	var cg color.Gray
	mono = true
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
	///

	////

	epdimg = waveshare.Mono2ByteImage(gimg)

	f, e := os.Create("output.bmp")
	glog.Errorln(e)
	bmp.Encode(f, gimg)
	f.Close()

	return epdimg

}
