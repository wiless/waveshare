package main

import (
	"image"
	"image/color"

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
	epd.DisplayFrame()
}

func Background() image.Gray {
	// img := epd.GetFrame()
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

	draw2dkit.Rectangle(gc, 0, 0, 60, 60)
	gc.FillStroke()

	gc.Save()

	gc.SetFontSize(20)
	gc.SetFillColor(color.Black)
	gc.SetStrokeColor(color.Black)

	gc.StrokeStringAt("40", 4, 10)
	gc.FillStroke()

	mimg := ws.ConvertToGray(img)
	bimg := ws.Mono2ByteImagev2(mimg)
	return bimg
}
