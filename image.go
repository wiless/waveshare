package waveshare

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

	img, e := jpeg.Decode(f)

	if e != nil {
		log.Println("Decode Error, trying png ", e)
	img, e = png.Decode(f)
if e!=nil{
		log.Panic("Decode Error ", e)
}
		return nil
	}

	log.Println("Bounds ", img.Bounds())
	
	//rorate it it...
	byteimg=image.NewGray(image.Rect(0,0,200,25))
	width, height := img.Bounds().Dx(), img.Bounds().Dy()

	
	bitcount := 0
	var bitstr [8]string
	
	bytecnt := 0
	var cg color.Gray
	if width >= 200 && height >= 200 {
		for row := 0; row < 200; row++ {
//			fmt.Println()
			bytecnt=0

			for col := 0; col < 200; col++ {
				c := img.At(col,row)
				r, _, _, _ := c.RGBA()
				if r > 0 {
					bitstr[bitcount] = "1"
				} else {
					bitstr[bitcount] = "0"
				}
				bitcount++
				if bitcount>7 {
					bitcount = 0
					str := strings.Join(bitstr[:], "")
					val, _ := strconv.ParseUint(str, 2, 8)
					cg.Y=uint8(val)
					byteimg.SetGray(row,bytecnt,cg)

					bytecnt++
					
				}				
				



			}

		}
	}

	
	return byteimg
}

