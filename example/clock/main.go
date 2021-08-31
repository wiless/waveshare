package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"strings"
	"time"

	"github.com/llgcode/draw2d"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"

	ws "github.com/wiless/waveshare"
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
	time.Sleep(2 * time.Second)
	epd.SetFrame(bimg)
	for {
		// time.Sleep(1 * time.Second)
		// updateTime()
		updateTimeBox()
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

	draw2dkit.Rectangle(gc, 0, 0, 130, 60)

	gc.FillStroke()
	gc.Save()

	// gc.SetFontSize(20)
	// gc.SetLineWidth(4)
	// gc.SetFillColor(color.Black)
	// gc.SetStrokeColor(color.Black)
	// gc.StrokeStringAt("2021", 15, 30)
	// gc.FillStroke()

	// Display Date & Day
	gc.SetFontSize(17)
	gc.SetLineWidth(3)
	gc.SetFillColor(color.Black)
	gc.SetStrokeColor(color.Black)
	now := time.Now()
	datestr := fmt.Sprintf("%s", now.Format("Jan 02"))
	gc.StrokeStringAt(datestr, 132, 20)
	gc.FillStroke()

	gc.SetFontSize(12)
	datestr = fmt.Sprintf("%s", now.Format("Monday"))

	gc.StrokeStringAt(datestr, 132, 45)
	gc.FillStroke()
	gc.Save()

	gc.SetStrokeColor(color.Black)
	gc.SetFillColor(color.White)
	gc.SetLineWidth(2)
	gc.MoveTo(0, 60)
	gc.LineTo(200, 60)
	gc.FillStroke()

	// /// Show Time
	// gc.SetFontSize(40)
	// // gc.SetLineWidth(7)
	// datestr = fmt.Sprintf("%s", now.Format("15:04:05"))
	// left, top, right, bottom := gc.GetStringBounds(datestr)
	// height := bottom - top
	// width := right - left
	// gc.SetStrokeColor(color.Black)
	// gc.SetFillColor(color.Black)

	// // log.Println("Locations left,top,right,bottom,height,width", left, top, right, bottom, height, width)
	// // log.Println("Position is ", 100-width/2, 80+70-height/2)
	// gc.StrokeStringAt(datestr, 100-width/2, 80+70-height/2)
	// gc.FillStroke()
	// gc.Save()

	// load cloud image
	// png.Decode()

	gc.SetFontSize(15)
	gc.SetLineWidth(3)
	msg := "This is good to render in two lines"

	if len(os.Args) > 1 {
		msg = os.Args[1]
	}
	wrapText(gc, msg, 0, 64, 200)

	gc.Close()
	mimg := ws.ConvertToGray(img)
	bimg := ws.Mono2ByteImagev2(mimg)
	return bimg
}

func updateTime() {

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

	// log.Println("Locations left,top,right,bottom,height,width", left, top, right, bottom, height, width)
	// log.Println("Position is ", 100-width/2, 70-height/2)
	gc.StrokeStringAt(datestr, 100-width/2, 70-height/2)
	gc.FillStroke()
	gc.Save()

	gimg := ws.ConvertToGray(updateImage)
	epd.SetSubFrame(60, 0, gimg)
}

func updateTimeBox() {

	epd.Init(false)
	updateImage := image.NewRGBA(image.Rect(0, 0, 130, 60))
	WIDTH := 130.0 / 2
	HEIGHT := 60.0 / 2
	gc := draw2dimg.NewGraphicContext(updateImage)
	gc.SetFont(ws.EPD_FONT)
	gc.SetFillColor(color.White)
	gc.SetStrokeColor(color.Black)

	gc.ClearRect(2, 2, 130, 60)
	gc.FillStroke()
	gc.Save()

	/// Show Time
	gc.SetFontSize(23)
	gc.SetLineWidth(4)
	now := time.Now()
	datestr := fmt.Sprintf("%s", now.Format("03:04:05")) // PM
	left, top, right, bottom := gc.GetStringBounds(datestr)
	height := bottom - top
	width := right - left
	gc.SetStrokeColor(color.Black)
	gc.SetFillColor(color.Black)

	log.Println("Locations left,top,right,bottom,height,width", left, top, right, bottom, height, width)
	log.Println("Position is ", WIDTH-width/2, HEIGHT-height/2)
	gc.StrokeStringAt(datestr, WIDTH-width/2, HEIGHT-height/20-top/2)
	// gc.StrokeStringAt(datestr, 2, 25)
	gc.FillStroke()
	gc.Save()

	gimg := ws.ConvertToGray(updateImage)
	epd.SetSubFrame(0, 0, gimg)
}

// wrap text
func wrapText(gc *draw2dimg.GraphicContext, text string, x, y, maxWidth float64) {
	words := strings.Split(text, " ")
	line := ""
	fmt.Println("Words ", words)
	var left, top, right, bottom float64
	for n := 0; n < len(words); n++ {
		testLine := line + words[n] + " "
		//   var metrics = gc.GetStringBounds(testLine);

		left, top, right, bottom = gc.GetStringBounds(testLine)
		height := bottom - top
		width := right - left
		if n == 0 {
			y = y - top
		}
		testWidth := width
		if testWidth > maxWidth && n > 0 {
			// gc.fillText(line, x, y);
			gc.StrokeStringAt(line, x, y)
			line = words[n] + " "
			y += height //lineHeight
		} else {
			line = testLine
		}
	}
	gc.StrokeStringAt(line, x, y)
	// context.fillText(line, x, y);
	gc.FillStroke()

}
