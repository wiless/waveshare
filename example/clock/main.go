package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/llgcode/draw2d"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"

	"github.com/wiless/waveshare"
)

var epd ws.EPD

func init() {
	draw2d.SetFontFolder(".")
}
func main() {
	ws.InitHW()
	epd.Init(true)

	bimg := Background()
	// ws.AsciiPrintByteImage("Background", bimg)
	epd.SetFrame(bimg)
	time.Sleep(1 * time.Second)
	epd.SetFrame(bimg)
	for {
		time.Sleep(1 * time.Second)
		Refresh()
		epd.DisplayFrame()
	}
}

func GetWeatherIcon() {

}

func Background() image.Gray {
	// img := epd.GetFrame()
	// cloud := ws.LoadImage("weathersmall.png")
	// if cloud == nil {
	// 	log.Panicln("Unable to load weather.png")
	// }
	// ws.AsciiPrintByteImage("cloud", *cloud)

	img := image.NewRGBA(image.Rect(0, 0, 200, 200))

	gc := draw2dimg.NewGraphicContext(img)
	gc.SetFont(ws.EPD_FONT)
	gc.SetFillColor(color.White)
	gc.ClearRect(0, 0, 200, 200)
	gc.FillStroke()
	gc.Save()
	gc.SetFillColor(color.White)
	gc.SetStrokeColor(color.Black)
	gc.SetLineWidth(2)

	draw2dkit.Rectangle(gc, 0, 0, 90, 60)

	gc.FillStroke()
	gc.Save()
	gc.SetFontSize(20)
	gc.SetLineWidth(4)
	gc.SetFillColor(color.Black)
	gc.SetStrokeColor(color.Black)
	gc.StrokeStringAt("40", 15, 30)
	gc.FillStroke()

	// Display Date & Day
	gc.SetFontSize(15)
	gc.SetLineWidth(3)
	gc.SetFillColor(color.Black)
	gc.SetStrokeColor(color.Black)
	now := time.Now()
	datestr := fmt.Sprintf("%s", now.Format("Jan 02"))
	gc.StrokeStringAt(datestr, 100, 20)
	gc.FillStroke()
	gc.SetFontSize(20)

	datestr = fmt.Sprintf("%s", now.Format("Monday"))
	gc.StrokeStringAt(datestr, 100, 45)
	gc.FillStroke()
	gc.Save()
	gc.SetStrokeColor(color.Black)

	gc.SetFillColor(color.White)
	gc.SetLineWidth(2)
	gc.MoveTo(0, 60)
	gc.LineTo(200, 60)
	gc.FillStroke()

	/// Show Time
	// gc.SetFontSize(40)
	// gc.SetLineWidth(7)
	// datestr = fmt.Sprintf("%s", now.Format("15:04:05"))
	// left, top, right, bottom := gc.GetStringBounds(datestr)
	// height := bottom - top
	// width := right - left
	// gc.SetStrokeColor(color.Black)
	// gc.SetFillColor(color.Black)

	// log.Println("Locations left,top,right,bottom,height,width", left, top, right, bottom, height, width)
	// log.Println("Position is ", 100-width/2, 80+70-height/2)
	// gc.StrokeStringAt(datestr, 100-width/2, 80+70-height/2)
	// gc.FillStroke()
	// gc.Save()

	// load cloud image
	// png.Decode()
	gc.Close()
	mimg := ws.ConvertToGray(img)
	bimg := ws.Mono2ByteImagev2(mimg)
	return bimg
}

func Refresh() {

	epd.Init(false)
	updateImage := image.NewRGBA(image.Rect(0, 0, 200, 140))

	gc := draw2dimg.NewGraphicContext(updateImage)
	gc.SetFont(ws.EPD_FONT)
	gc.SetFillColor(color.White)
	gc.ClearRect(0, 0, 200, 140)
	gc.FillStroke()
	gc.Save()

	/// Show Time
	gc.SetFontSize(40)
	gc.SetLineWidth(7)
	now := time.Now()
	datestr := fmt.Sprintf("%s", now.Format("15:04:05"))
	left, top, right, bottom := gc.GetStringBounds(datestr)
	height := bottom - top
	width := right - left
	gc.SetStrokeColor(color.Black)
	gc.SetFillColor(color.Black)

	log.Println("Locations left,top,right,bottom,height,width", left, top, right, bottom, height, width)
	log.Println("Position is ", 100-width/2, 70-height/2)
	gc.StrokeStringAt(datestr, 100-width/2, 70-height/2)
	gc.FillStroke()
	gc.Save()

	gimg := ws.ConvertToGray(updateImage)
	epd.SetSubFrame(60, 0, gimg)
}
