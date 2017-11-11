package main

import (
	"image"
	"image/color"
	"os"
	"time"

	"github.com/llgcode/draw2d"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"

	"log"

	"github.com/golang/freetype/truetype"
	"github.com/golang/glog"

	"golang.org/x/image/bmp"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/wiless/waveshare"
)

var mono = true
var epd ws.EPD

func main() {
	ws.InitHW()
	draw2d.SetFontFolder(".")
	epd.Init(true)
	epdimg := ImageGenerate()
	kavimg := ws.LoadImage("kavishbw.jpg")
	log.Println("Loading kavish..")
	ws.AsciiPrintByteImage("KAVISH", *kavimg)
	UpdateImage(epdimg)
	time.Sleep(2 * time.Second)

	log.Println("Loading Geometry.....")
	UpdateImage(epdimg)
	_ = kavimg
	epd.DisplayFrame()
	time.Sleep(2 * time.Second)
	// return
	//	 for {
	//	 	time.Sleep(5 * time.Second)
	//	 	log.Println("Toggling Image...")
	//	 	epd.DisplayFrame()
	//	 }

	_ = epdimg
	time.Sleep(2 * time.Second)
	for {

		PartialUpdate()
		time.Sleep(1 * time.Second)
	}

}
func UpdateImage(epdimg image.Gray) {

	epd.SetFrame(epdimg) // set both frames with same image
}

func PartialUpdate() {

	epd.Init(false)
	timeimg := image.NewRGBA(image.Rect(0, 0, 104, 50))
	gc := draw2dimg.NewGraphicContext(timeimg)
	gc.ClearRect(0, 0, 104, 50)
	gc.SetFillColor(color.Black)
	gc.SetStrokeColor(color.Black)
	draw2dkit.Rectangle(gc, 10, 10, 90, 40)
	gc.SetLineWidth(2)
	tstr := time.Now().Format("15:04:05 PM")
	gc.StrokeStringAt(tstr, 30, 30)
	gc.Stroke()
	gc.Save()
	draw2dimg.SaveToPngFile("subimage.png", timeimg)
	gimg := ws.ConvertToGray(timeimg)
	ws.SaveBMP("subimage.bmp", gimg)
	ws.AsciiPrint("Partial COLOR", timeimg)
	ws.AsciiPrint("Partial GRAY", gimg)
	epd.SetSubFrame(40, 8, gimg)
	epd.DisplayFrame()
}

func ImageGenerate() (epdimg image.Gray) {
	// img := image.NewGray(image.Rect(0, 0, 200, 210))
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	for r := 0; r < 200; r++ {
		for c := 0; c < 200; c++ {
			img.Set(c, r, color.White)
		}
	}

	gc := draw2dimg.NewGraphicContext(img)
	// gc.ClearRect(0, 0, 200, 200)
	// gc.Rotate(3.141)
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
	msg := " ABCDEFGHIJKLMNOP "
	// L, T, R, B := gc.GetStringBounds(msg)

	// fmt.Println("L T R B", L, T, R, B)
	gc.StrokeStringAt(msg, 0, 20)
	// gc.FillStroke()
	gc.SetFontSize(20)
	gc.SetLineWidth(4)
	datestr := time.Now().Format(time.Stamp)
	gc.StrokeStringAt(datestr, 10, 170)
	gc.FillStroke()
	gc.Close()

	ws.AsciiPrint("GEOMETRY ", img)

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
			oldPixel := img.At(c, r)

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
			gimg.SetGray(c, r, cg)

		}
	}
	///
	ws.AsciiPrint("GRAY GEOMETRY ", gimg)
	////

	epdimg = ws.Mono2ByteImage(gimg)
	// epdimg = ws.Mono2ByteImagev2(gimg)

	ws.AsciiPrintByteImage("BYTE EPDD ", epdimg)

	f, e := os.Create("output.bmp")
	glog.Errorln(e)
	bmp.Encode(f, gimg)
	f.Close()

	return epdimg

}
