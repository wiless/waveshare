package main

import (
	"image"
	"image/color"
	"os"

	"github.com/llgcode/draw2d"

	"github.com/lucasb-eyer/go-colorful"

	"github.com/llgcode/draw2d/draw2dimg"

	"github.com/golang/glog"

	"golang.org/x/image/bmp"

	"github.com/wiless/waveshare"
)

var mono = true
var epd waveshare.EPD

func main() {
	waveshare.InitHW()
	draw2d.SetFontFolder(".")

	epdimg := ImageGenerate()
	UpdateImage(epdimg)
}
func UpdateImage(epdimg image.Gray) {
	epd.Init(true)
	epd.ClearFrame(0xff)
	epd.SetFrame(epdimg)
	epd.DisplayFrame()
}

func PartialUpdate(img image.Gray, x, y uint8) {

	// 	  epd.init(epd.lut_partial_update)
	//     image = Image.open('monocolor.bmp')
	// ##
	//  # there are 2 memory areas embedded in the e-paper display
	//  # and once the display is refreshed, the memory area will be auto-toggled,
	//  # i.e. the next action of SetFrameMemory will set the other memory area
	//  # therefore you have to set the frame memory twice.
	//  ##
	//     epd.set_frame_memory(image, 0, 0)
	//     epd.display_frame()
	//     epd.set_frame_memory(image, 0, 0)
	//     epd.display_frame()

	//     time_image = Image.new('1', (96, 32), 255)  # 255: clear the frame
	//     draw = ImageDraw.Draw(time_image)
	//     font = ImageFont.truetype('/usr/share/fonts/truetype/freefont/FreeMonoBold.ttf', 32)
	//     image_width, image_height  = time_image.size
	//     while (True):
	//         # draw a rectangle to clear the image
	//         draw.rectangle((0, 0, image_width, image_height), fill = 255)
	//         draw.text((0, 0), time.strftime('%M:%S'), font = font, fill = 0)
	//         epd.set_frame_memory(time_image.rotate(90), 80, 80)
	//         epd.display_frame()
}

func ImageGenerate() (epdimg image.Gray) {
	// img := image.NewGray(image.Rect(0, 0, 200, 210))
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))

	gc := draw2dimg.NewGraphicContext(img)
	gc.SetFillColor(color.White)
	red, _ := colorful.Hex("#000000")
	gc.SetStrokeColor(red)
	gc.MoveTo(0, 0)
	gc.LineTo(150, 105)
	gc.QuadCurveTo(100, 20, 50, 20)
	gc.StrokeStringAt("Hello Sendil", 10, 10)

	gc.Close()
	gc.Stroke()
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
