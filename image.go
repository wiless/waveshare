//Package ws - Utility functions to handle image
package ws

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"golang.org/x/image/bmp"
	//	"github.com/kidoman/embd"
	//	"github.com/lucasb-eyer/go-colorful"
	//	"golang.org/x/image/bmp"
)

func LoadImage(ifname string) (byteimg *image.Gray) {

	f, e := os.Open(ifname)
	if e != nil {
		fmt.Println("File Read Error ", e)
		return nil
	}
	cfg, str, e := image.DecodeConfig(f)
	if e != nil {
		log.Println("Error Decoding file ", e, str, cfg)
		return nil
	}

	var img image.Image
	log.Println("Found ", cfg)
	switch str {
	case "jpeg":
		{
			f.Close()
			f, _ = os.Open(ifname)
			img, e = jpeg.Decode(f)
			if e != nil {
				log.Println("Decode Error, trying jpg ", e)
			}
		}
	case "png":
		{
			f.Close()
			f, _ = os.Open(ifname)
			img, e = png.Decode(f)
			if e != nil {
				log.Println("Decode Error, trying png ", e)
				return nil
			}
		}
	default:
		log.Println("Unknown Image type ", str)
		return nil
	}

	log.Println("Bounds ", img.Bounds())

	//rorate it it...
	byteimg = image.NewGray(image.Rect(0, 0, 25, 200))
	width, height := img.Bounds().Dx(), img.Bounds().Dy()

	bitcount := 0
	var bitstr [8]string

	bytecnt := 0
	var cg color.Gray
	if width >= 200 && height >= 200 {
		for row := 0; row < 200; row++ {
			//			fmt.Println()
			bytecnt = 0

			for col := 0; col < 200; col++ {
				c := img.At(col, row)

				r, g, b, a := c.RGBA()
				if r > 0 || g > 0 || b > 0 || a > 0 {
					bitstr[bitcount] = "1"
				} else {
					bitstr[bitcount] = "0"
				}
				bitcount++
				if bitcount > 7 {
					bitcount = 0
					str := strings.Join(bitstr[:], "")
					val, _ := strconv.ParseUint(str, 2, 8)
					cg.Y = uint8(val)
					byteimg.SetGray(bytecnt, row, cg)

					bytecnt++

				}

			}

		}
	}

	return byteimg
}

func ConvertToGray(img image.Image) *image.Gray {
	b := img.Bounds()
	gimg := image.NewGray(b)
	var cg color.Gray
	var mono = true
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

func AsciiPrint(name string, img image.Image) {
	b := img.Bounds()
	R, C := b.Max.Y, b.Max.X

	fmt.Printf("\n %s = [rows x cols] = %d,%d \n", name, R, C)
	for r := 0; r < R; r++ {
		fmt.Printf("\n Row %03d : ", r)
		for c := 0; c < C; c++ {
			clr := img.At(c, r)
			pix, _, _, _ := clr.RGBA()
			if pix > 0 {
				pix = 1
			}
			fmt.Printf("%d", pix)
		}
	}
}

func AsciiPrintByteImage(name string, img image.Gray) {
	b := img.Bounds()
	R, C := b.Max.Y, b.Max.X
	fmt.Printf("\n %s = [rows x cols] = %d,%d \n", name, R, C)
	for r := 0; r < R; r++ {
		fmt.Printf("\n Row %03d : ", r)
		for c := 0; c < C; c++ {
			clr := img.GrayAt(c, r).Y
			fmt.Printf("%08b", clr)
		}
	}
}
