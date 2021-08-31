package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"strings"
	"time"

	"github.com/llgcode/draw2d/draw2dimg"
	ws "github.com/wiless/waveshare"
)

// func LoadEPDBin(fname string) image.Gray {

// 	var img *image.Gray
// 	img = image.NewGray(image.Rect(0, 0, 25, 200))
// 	b := img.Bounds()
// 	R, C := b.Max.Y, b.Max.X
// 	fmt.Printf("\n %s = [rows x cols] = %d,%d \n", fname, R, C)
// 	var rowbytes []byte
// 	rowbytes, fer := os.ReadFile(fname)
// 	Err(fer)
// 	cnt := 0
// 	for r := 0; r < R; r++ {
// 		fmt.Printf("\n Row %03d : ", r)

// 		for c := 0; c < C; c++ {

// 			var clr color.Gray
// 			clr.Y = uint8(rowbytes[cnt])
// 			// clr := img.GrayAt(c, r).Y
// 			img.SetGray(c, r, clr)
// 			fmt.Printf("%08b", clr.Y)
// 			cnt++
// 		}

// 	}
// 	return *img
// }

func LoadEPDBin(rawbytes []byte) image.Gray {

	var img *image.Gray
	img = image.NewGray(image.Rect(0, 0, 25, 200))
	b := img.Bounds()
	R, C := b.Max.Y, b.Max.X
	fmt.Printf("\n = [rows x cols] = %d,%d \n", R, C)

	cnt := 0
	for r := 0; r < R; r++ {
		fmt.Printf("\n Row %03d : ", r)

		for c := 0; c < C; c++ {

			var clr color.Gray
			clr.Y = uint8(rawbytes[cnt])
			// clr := img.GrayAt(c, r).Y
			img.SetGray(c, r, clr)
			fmt.Printf("%08b", clr.Y)
			cnt++
		}

	}
	return *img
}

func updateTimeBox(interval int) {
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
