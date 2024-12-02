package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/llgcode/draw2d"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"

	"github.com/golang/freetype/truetype"
	"github.com/golang/glog"

	"golang.org/x/image/bmp"
	"golang.org/x/image/font/gofont/goregular"

	ws "github.com/wiless/waveshare"
)

var mono = true
var epd ws.EPD
var Err = func(e error) {
	if e != nil {
		fmt.Println("Error %v", e)
	}

}

func LoadEPDBin(fname string) image.Gray {

	var img *image.Gray
	img = image.NewGray(image.Rect(0, 0, 25, 200))
	b := img.Bounds()
	R, C := b.Max.Y, b.Max.X
	fmt.Printf("\n %s = [rows x cols] = %d,%d \n", fname, R, C)
	var rowbytes []byte
	rowbytes, fer := os.ReadFile(fname)
	Err(fer)
	cnt := 0
	for r := 0; r < R; r++ {
		fmt.Printf("\n Row %03d : ", r)

		for c := 0; c < C; c++ {

			var clr color.Gray
			clr.Y = uint8(rowbytes[cnt])
			// clr := img.GrayAt(c, r).Y
			img.SetGray(c, r, clr)
			fmt.Printf("%08b", clr.Y)
			cnt++
		}

	}
	return *img
}

var imageframes []image.Gray

func main() {
	ws.InitHW()
	draw2d.SetFontFolder(".")
	epd.Init(true)

	if len(os.Args) > 1 {
		if len(os.Args) == 2 {
			bimg := LoadEPDBin(os.Args[1])
			imageframes = append(imageframes, bimg)
			epd.SetFrame(bimg)
			// UpdateImage(imageframes[0])
			time.Sleep(2 * time.Second)
			epd.SetFrame(bimg)
			// UpdateImage(imageframes[0])

		} else {
			for i := 1; i < len(os.Args); i++ {
				tmpimg := LoadEPDBin(os.Args[i])
				imageframes = append(imageframes, tmpimg)
				UpdateImage(tmpimg)
				time.Sleep(1 * time.Second)
			}
		}

	} else {
		epdimg := ImageGenerate()
		imageframes = append(imageframes, epdimg)
		UpdateImage(epdimg)
		// NOT updated to DISPLAY buffer !!
	}

	// IF only one image.. upate time on partial area
	if len(imageframes) == 1 {
		for {
			time.Sleep(1000 * time.Millisecond)
			PartialUpdate()
			// epd.DisplayFrame()
			// epd.DisplayFrame()
		}

	}
	// IF only two images.. keep switching between them from IMAGEs in the memory
	if len(imageframes) == 2 {

		for {
			time.Sleep(5 * time.Second)
			log.Println("Toggling Image...")
			epd.DisplayFrame()
		}
	}
	// IF MORE than two images.. keep CYCLING loading images from imageframes
	if len(imageframes) > 2 {
		cnt := 0
		for {
			time.Sleep(3 * time.Second)
			log.Println("Next Image...")
			UpdateImage(imageframes[cnt])
			cnt++
			if cnt == len(imageframes) {
				cnt = 0
			}
		}
	}

}
func UpdateImage(epdimg image.Gray) {

	epd.SetFrame(epdimg) // set both frames with same image
}

func PartialUpdate() {
	epd.Init(false)
	timeimg := image.NewRGBA(image.Rect(0, 0, 104, 55))
	gc := draw2dimg.NewGraphicContext(timeimg)
	gc.ClearRect(0, 0, 104, 55)
	gc.SetFillColor(color.Black)
	gc.SetStrokeColor(color.Black)
	draw2dkit.Rectangle(gc, 2, 2, 102, 53)
	gc.SetLineWidth(2)
	tstr := time.Now().Format("15:04:05 PM")
	gc.SetFontSize(13)
	gc.StrokeStringAt(tstr, 13, 35)
	gc.Stroke()
	gc.Save()
	// draw2dimg.SaveToPngFile("subimage.png", timeimg)
	gimg := ws.ConvertToGray(timeimg)
	// ws.SaveBMP("subimage.bmp", gimg)
	// ws.AsciiPrint("Partial COLOR", timeimg)
	// ws.AsciiPrint("Partial GRAY", gimg)
	epd.SetSubFrame(0, 0, gimg)
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

	// ws.AsciiPrint("GEOMETRY ", img)

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
	// ws.AsciiPrint("GRAY GEOMETRY ", gimg)
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
