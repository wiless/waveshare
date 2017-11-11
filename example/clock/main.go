package main

import (
	"image"
	"image/color"

	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"

	"github.com/wiless/waveshare"
)

var epd ws.EPD

func init() {

}
func main() {
	ws.InitHW()
	epd.Init(true)
	bimg := Background()
	epd.SetFrame(bimg)
	ws.AsciiPrintByteImage("Background", bimg)
	epd.DisplayFrame()
}

func Background() image.Gray {
	// img := epd.GetFrame()
	img := image.NewRGBA(image.Rect(0, 0, 60, 60))
	gc := draw2dimg.NewGraphicContext(img)
	gc.SetFillColor(color.White)
	gc.SetStrokeColor(color.Black)
	draw2dkit.Rectangle(gc, 0, 0, 58, 59)
	gc.FillStroke()
	gc.Close()
	mimg := ws.ConvertToGray(img)
	bimg := ws.Mono2ByteImagev2(mimg)
	return bimg
}
