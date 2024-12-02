package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/lucasb-eyer/go-colorful"

	"golang.org/x/image/bmp"
)

const (
	channel = 0
	speed   = 2000000
	bpw     = 8
	delay   = 0
)

func main() {
	fmt.Printf("bmp tools")
	//	LoadImage("foodklubbw.jpg", "IMAGE_DATA")

}

func LoadImage(ifname, ofname string) []byte {
	of, _ := os.Create(ofname + ".dat")

	f, e := os.Open(ifname)
	if e != nil {
		fmt.Println("File Read Error ", e)
		return nil
	}

	img, e := jpeg.Decode(f)
	if e != nil {
		log.Println("Decode Error ", e)
		return []byte{}
	}
	log.Println("Bounds ", img.Bounds())
	width, height := img.Bounds().Dx(), img.Bounds().Dy()

	// xx,_:=strconv.ParseUint("10110001", 2, 8)
	bitcount := 0
	var bitstr [8]string
	fmt.Fprintf(of, "{")
	var result []byte
	result = make([]byte, 200*200/8)
	bytecnt := 0

	if width >= 200 && height >= 200 {
		for row := 0; row < 200; row++ {
			fmt.Println()
			for col := 0; col < 200; col++ {
				c := img.At(200-col, 200-row)
				r, g, b, _ := c.RGBA()
				_ = g
				_ = b
				// fmt.Printf("[%v %v %v]", r, g, b)
				if r > 0 {
					bitstr[bitcount] = "1"
					// fmt.Print(0)
				} else {
					bitstr[bitcount] = "0"
					// fmt.Print(1)
				}
				bitcount++
				if bitcount == 8 {
					bitcount = 0
					str := strings.Join(bitstr[:], "")
					// log.Println(str)
					val, _ := strconv.ParseUint(str, 2, 8)
					result[bytecnt] = byte(val)
					bytecnt++
					fmt.Printf("%08b", byte(val))
					fmt.Fprintf(of, "0x%02x,", val)
				}

			}
			// fmt.Println()
			fmt.Fprintf(of, "\n")

		}
	}

	CreateCPP("imagedata.cpp", ofname, result)
	fmt.Fprintf(of, "};")

	return []byte{}
}

func CreateCPP(fname string, varname string, val []byte) {
	f, e := os.Create(fname)
	if e != nil {
		log.Printf("Error Saving ", e)
		return
	}
	header := `
	/**
	File Autogenerated by bmptest.go
	DO NOT EDIT
	**/
	#include "imagedata.h"
	`
	fmt.Fprintf(f, header)
	fmt.Fprintf(f, "\n const unsigned char IMAGE_DATA[] = {")
	count := len(val)

	for i, v := range val {
		if i < count-1 {
			fmt.Fprintf(f, "%v,", v)
		} else {
			fmt.Fprintf(f, "%v", v)
		}
	}

	fmt.Fprintf(f, " };")
}

func SaveArray(fname string) {
	f, e := os.Create(fname)
	if e != nil {
		log.Println(e)
		return
	}

	// var r image.Rectangle
	r := image.Rect(0, 0, 200, 200)
	log.Println("Length of IMGdata ", len(IMAGE_DATA))
	img := image.NewGray(r)

	fmt.Printf("\n input = %x, output =%s", IMAGE_DATA[0:1], strconv.FormatInt(int64(IMAGE_DATA[0]), 2))
	nextbyte := 0
	var bitarray []string
	var bitptr int

	for x := 0; x < 200; x++ {
		fmt.Println("%03d   ", x)
		for y := 0; y < 200; y++ {
			if bitptr == 0 {

				temp := fmt.Sprintf("%08b", IMAGE_DATA[nextbyte])
				bitarray = strings.Split(temp, "")
				nextbyte++
			}
			var clr color.Color

			if bitarray[bitptr] == "0" {
				fmt.Printf("0")

				clr = colorful.LinearRgb(1, 1, 1)
			} else {
				fmt.Printf("1")

				clr = colorful.LinearRgb(0, 0, 0)
			}
			img.Set(x, y, clr)

			// } else {
			// 	clr := colorful.LinearRgb(rand.Float64(), rand.Float64(), rand.Float64())
			// 	img.Set(x, y, clr)
			// }

			bitptr++
			if bitptr > 7 {
				bitptr = 0
			}
		}
	}
	fmt.Println()
	e = bmp.Encode(f, img)
	if e != nil {
		fmt.Print(e)
	}

	f.Close()
}

// /* 0x00,0x01,0xC8,0x00,0xC8,0x00, */
var IMAGE_DATA []byte = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFE, 0x01, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xF8, 0x00, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF7, 0xFF,
	0x3F, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xF0, 0x00,
	0x3F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xE1, 0xFE, 0x1F, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xF0, 0x00, 0x1F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xE3, 0xFF, 0x1F, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1,
	0xE0, 0x00, 0x1F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xE3, 0xFF, 0x9F, 0xFF, 0xFF, 0xFF, 0xFF,
	0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xE1, 0xE3, 0x00, 0x00, 0x00, 0x00, 0x01,
	0xFF, 0xFF, 0xE7, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xF1, 0xC3, 0xF3, 0x0F, 0xFF, 0xFF, 0xFF, 0xFC, 0xFF, 0xFF, 0xE3, 0xFF, 0x8F, 0xFF, 0xFF,
	0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xC3, 0x3B, 0x0F, 0xFF, 0xFF,
	0xFF, 0xFE, 0xFF, 0xFF, 0xE3, 0xFF, 0x1F, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xF1, 0xC3, 0x0F, 0x0F, 0xFF, 0xFF, 0xFF, 0xFE, 0xFF, 0xFF, 0xE1, 0xFE, 0x1F,
	0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xC3, 0x0F, 0x0F,
	0xFF, 0xFF, 0xFF, 0xFE, 0xFF, 0xFF, 0xF0, 0xFC, 0x1F, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xE3,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xC3, 0x87, 0x0F, 0xFF, 0xFF, 0xFF, 0xFE, 0xFF, 0xFF, 0xF0,
	0x00, 0x3F, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xE0, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xC1,
	0x03, 0x0F, 0xFF, 0xFF, 0xFF, 0xFE, 0xFF, 0xFF, 0xF8, 0x00, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F,
	0xFF, 0xE0, 0x0F, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xE0, 0x00, 0x1F, 0xFF, 0xFF, 0xFF, 0xFE, 0xFF,
	0xFF, 0xFE, 0x01, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFC, 0x03, 0xFF, 0xFF, 0xFF, 0xFF,
	0xF1, 0xE0, 0x00, 0x1F, 0xFF, 0xFF, 0xFF, 0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0x8F, 0xFF, 0xFF, 0x80, 0x7F, 0xFF, 0xFF, 0xFF, 0xF1, 0xF0, 0x00, 0x3F, 0xFF, 0xFF, 0xFF,
	0x7E, 0xFF, 0xFF, 0xFB, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xF0, 0x1F, 0xFF,
	0xFF, 0xFF, 0xF1, 0xF8, 0x00, 0x7F, 0xFF, 0xFF, 0xFC, 0x3E, 0xFF, 0xFF, 0xE0, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0x00, 0x1F, 0xFF, 0xFF, 0xFF, 0xF1, 0xFE, 0x00, 0xFF, 0xFF,
	0xFF, 0xF0, 0x3E, 0xFF, 0xFF, 0xC6, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xF8, 0x07,
	0x8F, 0xFF, 0xFF, 0xFF, 0xF1, 0xFE, 0x87, 0xFF, 0xFF, 0xFF, 0xE0, 0xCE, 0xFF, 0xFF, 0xCE, 0x7F,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xE0, 0x1F, 0x8F, 0xC0, 0x07, 0xFF, 0xF1, 0xFE, 0xFF,
	0xFF, 0xFF, 0xFF, 0x81, 0x86, 0xFF, 0xFF, 0xCE, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF,
	0xE0, 0xFF, 0xCF, 0x80, 0x07, 0xFF, 0xF1, 0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0x02, 0x06, 0xFF, 0xFF,
	0xC6, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xE3, 0xFF, 0xEF, 0x00, 0x07, 0xFF, 0xF1,
	0xFE, 0xFF, 0xFF, 0xFF, 0xFE, 0x0C, 0x02, 0xFF, 0xFF, 0xE0, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0x1F, 0xFF, 0xFF, 0xF1, 0xFE, 0xFF, 0xFF, 0xFF, 0xFE, 0x18, 0x02,
	0xFF, 0xFF, 0xFB, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0x1F, 0xFF,
	0xFF, 0xF1, 0xFE, 0xFF, 0xFF, 0xFF, 0xFE, 0x20, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0x8F, 0xFF, 0xFE, 0x01, 0xFF, 0x1F, 0xFF, 0xFF, 0xF1, 0xFE, 0xFF, 0xFF, 0xFF, 0xFE,
	0xC0, 0x02, 0xFF, 0xFF, 0xFF, 0x03, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xF0, 0x01, 0xFF,
	0x9F, 0xFF, 0xFF, 0xF1, 0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x02, 0xFF, 0xFF, 0xF8, 0x00, 0x7F,
	0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xE0, 0x01, 0xE0, 0x80, 0x0F, 0xFF, 0xF1, 0xFE, 0xFF, 0xFF,
	0xFF, 0xFF, 0x00, 0x02, 0xFF, 0xFF, 0xF0, 0x00, 0x3F, 0xFF, 0xE7, 0xFF, 0xFF, 0x8F, 0xFF, 0xE0,
	0x03, 0xE0, 0x00, 0x07, 0xFF, 0xF1, 0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0x06, 0xFF, 0xFF, 0xE0,
	0xCC, 0x1F, 0xFF, 0xE0, 0xFF, 0xFF, 0x8F, 0xFF, 0xE2, 0x71, 0xE0, 0x00, 0x07, 0xFF, 0xF1, 0xFE,
	0xFF, 0xFF, 0xFF, 0xFF, 0xC0, 0x0E, 0xFF, 0xFF, 0xE3, 0xC7, 0x1F, 0xFF, 0xE0, 0x3F, 0xFF, 0x8F,
	0xFF, 0xE7, 0x39, 0xE0, 0x00, 0x0F, 0xFF, 0xF1, 0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xC0, 0x0E, 0xFF,
	0xFF, 0xE7, 0xE7, 0x8F, 0xFF, 0xF0, 0x07, 0xFF, 0x8F, 0xFF, 0xE7, 0x38, 0xFF, 0xFF, 0xFF, 0xFF,
	0xF1, 0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x1E, 0xFF, 0xFF, 0xE7, 0xE7, 0x8F, 0xFF, 0xFE, 0x01,
	0xFF, 0x8F, 0xFF, 0xE3, 0x10, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFE, 0xFF, 0xFF, 0xFF, 0xFC, 0x00,
	0x3E, 0xFF, 0xFF, 0xE3, 0xC7, 0x8F, 0xFF, 0xFF, 0xC0, 0x7F, 0x8F, 0xFF, 0xE3, 0x01, 0xFF, 0x1F,
	0xC7, 0xFF, 0xF1, 0xFE, 0xFF, 0x83, 0xFF, 0xF0, 0x00, 0x7E, 0xFF, 0xFF, 0xE0, 0x07, 0x1F, 0xFF,
	0xFF, 0xF8, 0x3F, 0x8F, 0xFF, 0xF3, 0x81, 0xFF, 0x1F, 0xC7, 0xFF, 0xF1, 0xFE, 0xFC, 0x01, 0xFF,
	0xE0, 0x00, 0xFE, 0xFF, 0xFF, 0xF0, 0x0F, 0x1F, 0xFF, 0xFF, 0xC0, 0x1F, 0x8F, 0xFF, 0xFF, 0xC7,
	0xFF, 0x1F, 0xC7, 0xFF, 0xF1, 0xFE, 0xF0, 0x07, 0xFF, 0x80, 0x01, 0xFE, 0xFF, 0xFF, 0xF8, 0x1F,
	0xBF, 0xFF, 0xFE, 0x00, 0x0F, 0x8F, 0xFF, 0xFF, 0xFF, 0xF8, 0x00, 0x07, 0xFF, 0xF1, 0xFE, 0xC0,
	0x18, 0xFE, 0x00, 0x03, 0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF0, 0x07, 0x87, 0x8F, 0xFF,
	0xFF, 0xFF, 0xF8, 0x00, 0x0F, 0xFF, 0xF1, 0xFE, 0xC0, 0xE0, 0xF8, 0x00, 0x0F, 0xFE, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xE0, 0x1F, 0xC7, 0x8F, 0xFF, 0xFF, 0xFF, 0xF8, 0x00, 0x0F, 0xFF, 0xF1,
	0xFE, 0xC7, 0x00, 0xE0, 0x00, 0x1F, 0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0x9F, 0xFF, 0xE0, 0xFF, 0xC7,
	0x8F, 0xFE, 0x00, 0x01, 0xFF, 0x1E, 0xFF, 0xFF, 0xF1, 0xFE, 0xDC, 0x00, 0xC0, 0x00, 0x3F, 0xFE,
	0xFF, 0xFF, 0xF8, 0x3F, 0x1F, 0xFF, 0xE3, 0xFF, 0xC7, 0x8F, 0xFE, 0x00, 0x01, 0xFF, 0x1F, 0xFF,
	0xFF, 0xF1, 0xFE, 0xE0, 0x00, 0x00, 0x00, 0xFF, 0xFE, 0xFF, 0xFF, 0xF0, 0x0F, 0x1F, 0xFF, 0xFF,
	0xFF, 0xC7, 0x8F, 0xFE, 0x00, 0x01, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFE, 0xC0, 0x00, 0x00, 0x03,
	0xFF, 0xFE, 0xFF, 0xFF, 0xE0, 0x03, 0x1F, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFE, 0x03, 0xFF, 0xF8,
	0x7F, 0xFF, 0xFF, 0xF1, 0xFE, 0xE0, 0x00, 0x00, 0x07, 0xFF, 0xFE, 0xFF, 0xFF, 0xE3, 0xC1, 0x9F,
	0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xC0, 0x7F, 0xF8, 0x1F, 0xFF, 0xFF, 0xF1, 0xFE, 0xE0, 0x00,
	0x00, 0x1F, 0xFF, 0xFE, 0xFF, 0xFF, 0xE7, 0xE0, 0x9F, 0xFF, 0xF8, 0x00, 0xFF, 0x8F, 0xFF, 0xF8,
	0x1F, 0xF8, 0x07, 0xFF, 0xFF, 0xF1, 0xFE, 0xF0, 0x00, 0x00, 0xFF, 0xFF, 0xFE, 0xFF, 0xFF, 0xE7,
	0xF8, 0x1F, 0xFF, 0xF0, 0x00, 0xFF, 0x8F, 0xFF, 0xFE, 0x03, 0xF8, 0x01, 0xFF, 0xFF, 0xF1, 0xFE,
	0xF8, 0x00, 0x03, 0xFF, 0xFF, 0xFE, 0xFF, 0xFF, 0xE7, 0xFC, 0x1F, 0xFF, 0xE0, 0x00, 0xFF, 0x8F,
	0xFF, 0xFF, 0x83, 0xF8, 0xC0, 0x7F, 0xFF, 0xF1, 0xFE, 0xFC, 0x00, 0x3F, 0xFF, 0xFF, 0xFE, 0xFF,
	0xFF, 0xE3, 0xFE, 0x1F, 0xFF, 0xE0, 0x00, 0xFF, 0x8F, 0xFF, 0xFC, 0x07, 0xF8, 0xF0, 0x1F, 0xFF,
	0xF1, 0xFE, 0xFF, 0xB7, 0xFF, 0xFF, 0xFF, 0xFE, 0xFF, 0xFF, 0xF1, 0xFF, 0x1F, 0xFF, 0xE3, 0x31,
	0xFF, 0x8F, 0xFF, 0xE0, 0x3F, 0xF8, 0xFC, 0x07, 0xFF, 0xF1, 0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFE, 0xFF, 0xFF, 0xF3, 0xFF, 0x9F, 0xFF, 0xE3, 0x38, 0xFF, 0x8F, 0xFF, 0x00, 0xFF, 0xF8, 0xFF,
	0x03, 0xFF, 0xF1, 0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xE3, 0x38, 0xFF, 0x8F, 0xFE, 0x07, 0xFF, 0xF8, 0xFF, 0xC3, 0xFF, 0xF1, 0xFE, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xE3, 0x18, 0xFF, 0x8F, 0xFE, 0x00, 0x01,
	0xF8, 0xFF, 0xF3, 0xFF, 0xF1, 0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE, 0xFF, 0xFF, 0xFF, 0xDF,
	0xFF, 0xFF, 0xE3, 0x00, 0xFF, 0x8F, 0xFE, 0x00, 0x01, 0xF8, 0xFF, 0xFF, 0xFF, 0xF1, 0xFE, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFE, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xE3, 0x80, 0xFF, 0x8F, 0xFE,
	0x00, 0x01, 0xF8, 0xFF, 0xFF, 0xFF, 0xF1, 0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE, 0xFF, 0xFF,
	0xFF, 0x87, 0xFF, 0xFF, 0xF1, 0x81, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1,
	0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE, 0xFF, 0xFF, 0xFF, 0xC3, 0xFF, 0xFF, 0xFF, 0xC3, 0xFF,
	0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFE, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE,
	0xFF, 0xFF, 0xFF, 0xE3, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xF1, 0xFE, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFC, 0xFF, 0xFF, 0xFF, 0xF3, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x01, 0xFF, 0xFF, 0xFF, 0xE3, 0xFF, 0xFC, 0x00, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xE3, 0xFF,
	0xFC, 0x00, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xC3, 0xFF, 0xFC, 0x00, 0x00, 0xFF, 0x8F, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0x87, 0xFF, 0xFC, 0x00, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xE3, 0xF1, 0xFF, 0x8F,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0x1F, 0xFF, 0xFF, 0xE3, 0xF9, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x1F, 0xFF, 0xFF, 0xE3, 0xF8,
	0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0x3F, 0xFF, 0xFF, 0xE3, 0xF8, 0xFF, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x01, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE, 0x3F, 0xFF, 0xFF,
	0xE1, 0xF0, 0xFF, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x1F, 0xFF, 0xFF, 0xF0, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x1F,
	0xFF, 0xFF, 0xF0, 0x01, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x0F, 0xFF, 0xFF, 0xF8, 0x01, 0xFF, 0x8F, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0x87, 0xFF, 0xFF, 0xFE, 0x07, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xCF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF8,
	0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF8, 0xFF, 0xFF, 0xF0, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF8, 0x70, 0x3F,
	0xFF, 0xE0, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFA, 0xFF, 0xFF, 0xFF, 0xFF, 0xF0, 0x20, 0x1F, 0xFF, 0xE0, 0x00, 0xFF, 0x8F, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xC0, 0x1F, 0xFF, 0xFF, 0xFF, 0xE0,
	0x02, 0x1F, 0xFF, 0xE3, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF,
	0xFF, 0xFF, 0xFF, 0x00, 0x0F, 0xFF, 0xFF, 0xFF, 0xE3, 0x87, 0x1F, 0xFF, 0xE3, 0xFF, 0xFF, 0x8F,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFE, 0x00, 0x07, 0xFF, 0xFF,
	0xFF, 0xE7, 0x8F, 0x8F, 0xFF, 0xE3, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xF1, 0xFF, 0xFF, 0xFF, 0xFE, 0x00, 0x02, 0x7F, 0xFF, 0xFF, 0xE7, 0xCF, 0x8F, 0xFF, 0xF0, 0xFF,
	0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFC, 0x00, 0x00,
	0x7F, 0xFF, 0xFF, 0xE7, 0xCF, 0x9F, 0xFF, 0xE0, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xF8, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xE3, 0xFF, 0x1F, 0xFF,
	0xE0, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xF8,
	0x30, 0x00, 0xFF, 0xFF, 0xFF, 0xF3, 0xFE, 0x1F, 0xFF, 0xE0, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xF8, 0x30, 0x00, 0xFF, 0xFF, 0xFF, 0xF7, 0xFF,
	0x3F, 0xFF, 0xE0, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF,
	0xFF, 0xF0, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xF0, 0x00, 0x00, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1,
	0xFF, 0xFF, 0xFF, 0x90, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x1F, 0xFF, 0xFF, 0xFF, 0xFF,
	0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0x10, 0x00, 0x00, 0xFF,
	0xFF, 0xFF, 0xF8, 0x1F, 0x1F, 0xFF, 0xE0, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xF1, 0xFF, 0xFF, 0xFC, 0x10, 0x00, 0x00, 0xFF, 0xFF, 0xFF, 0xF0, 0x0F, 0x1F, 0xFF, 0xE0,
	0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xF8, 0x18, 0x30,
	0x00, 0xFF, 0xFF, 0xFF, 0xE0, 0x03, 0x1F, 0xFF, 0xE0, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xF8, 0x08, 0x38, 0x00, 0xFF, 0xFF, 0xFF, 0xE3, 0xC1, 0x9F,
	0xFF, 0xE0, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xF0,
	0x08, 0x30, 0x01, 0xFF, 0xFF, 0xFF, 0xE7, 0xF0, 0x9F, 0xFF, 0xFF, 0xE1, 0xFF, 0x8F, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xE0, 0x08, 0x00, 0x01, 0xFF, 0xFF, 0xFF, 0xE7,
	0xF8, 0x1F, 0xFF, 0xFF, 0xF9, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF,
	0xFF, 0xE0, 0xE4, 0x00, 0x01, 0xFF, 0xFF, 0xFF, 0xE7, 0xFC, 0x1F, 0xFF, 0xFF, 0xF8, 0xFF, 0x8F,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xE0, 0xE2, 0x00, 0x03, 0xFF, 0xFF,
	0xFF, 0xE3, 0xFE, 0x1F, 0xFF, 0xFF, 0xF8, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xF1, 0xFF, 0xFF, 0xC0, 0xE2, 0x00, 0x07, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0x1F, 0xFF, 0xE0, 0x00,
	0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xC0, 0x41, 0x00, 0x0F,
	0xFF, 0xFF, 0xFF, 0xF3, 0xFF, 0x9F, 0xFF, 0xE0, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xC0, 0x00, 0xC0, 0x1F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xE0, 0x01, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xC0, 0x00,
	0x30, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xE0, 0x03, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xC0, 0x00, 0x07, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF,
	0xC0, 0x00, 0x00, 0x3F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xC0, 0x00, 0x00, 0x7F, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1,
	0xFF, 0xFF, 0xC0, 0x40, 0x00, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0xBF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xC0, 0xE0, 0x00, 0x7F, 0xFF,
	0xFF, 0xFF, 0xFF, 0x3F, 0xFF, 0xFF, 0xFF, 0x87, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xF1, 0xFF, 0xFF, 0xE0, 0xE0, 0x00, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0x3F, 0xFF, 0xFF, 0x1E,
	0x03, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xE0, 0xE0, 0x00,
	0x7F, 0xFF, 0xFF, 0xFE, 0xFF, 0xBF, 0xCF, 0xFE, 0x3E, 0x01, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xE0, 0x00, 0x00, 0x3F, 0xFF, 0xFF, 0xFE, 0x3F, 0xFF, 0x8F,
	0xFE, 0x3C, 0x01, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xF0,
	0x00, 0x00, 0x3F, 0xFF, 0xFF, 0xFF, 0x3C, 0x0F, 0x9F, 0xFE, 0x3C, 0x30, 0xFF, 0x8F, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xF8, 0x00, 0x01, 0x1F, 0xFF, 0xFF, 0xFF, 0xF0,
	0x03, 0xFF, 0xFE, 0x3C, 0x78, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF,
	0xFF, 0xF8, 0x00, 0x03, 0xFF, 0xFF, 0xFF, 0xFF, 0xE0, 0x00, 0xFF, 0xFE, 0x38, 0x78, 0xFF, 0x8F,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFC, 0x00, 0x07, 0xFF, 0xFF, 0xFF,
	0xFF, 0xC0, 0x00, 0xFF, 0xFE, 0x18, 0x78, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xF1, 0xFF, 0xFF, 0xFF, 0x00, 0x0F, 0xFF, 0xFF, 0xFF, 0xFF, 0xC0, 0x00, 0x7F, 0xFE, 0x00, 0xF8,
	0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0x80, 0x3F, 0xFF,
	0xFF, 0xFF, 0xFF, 0x80, 0x00, 0x3F, 0xFF, 0x00, 0xF8, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0x00, 0x3F, 0xFF,
	0x01, 0xF0, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0x00, 0x3F, 0xFF, 0xC3, 0xF1, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0x00, 0x00,
	0x31, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0x00, 0x00, 0x31, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0x80, 0x00, 0x3F, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x80, 0x00, 0x3F, 0xFF, 0xFF, 0xFF, 0xFF,
	0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0x80, 0x00, 0x7F, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xC0, 0x00, 0x7F, 0xFF, 0xFF,
	0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xE0, 0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF0, 0x01, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x78, 0x03, 0xDF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE, 0x3F,
	0x1F, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE, 0x7F, 0xFF, 0xCF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFE, 0xFF, 0x3F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x3F, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0x3F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0x00, 0x00,
	0x00, 0x01, 0x7F, 0xFF, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFD, 0x80, 0x00, 0x00, 0x03, 0x1F, 0xFF, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x0F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF8,
	0x80, 0x00, 0x00, 0x02, 0x0F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x40, 0x00, 0x00, 0x04, 0x0F, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1,
	0xFF, 0xF0, 0x60, 0x00, 0x00, 0x0C, 0x07, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x20, 0x00, 0x00, 0x08, 0x07,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xF1, 0xFF, 0xF0, 0x30, 0x00, 0x00, 0x10, 0x07, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x10, 0x00, 0x00,
	0x30, 0x07, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x18, 0x00, 0x00, 0x20, 0x07, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x08,
	0x00, 0x00, 0x60, 0x07, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x04, 0x00, 0x00, 0x40, 0x07, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF,
	0xF0, 0x06, 0x00, 0x00, 0x80, 0x07, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x02, 0x00, 0x01, 0x80, 0x07, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xF1, 0xFF, 0xF0, 0x03, 0x00, 0x01, 0x00, 0x07, 0xFF, 0xFF, 0xFF, 0xFF, 0xE0, 0x00, 0x00, 0x00,
	0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x01, 0x00, 0x02, 0x00,
	0x07, 0xFF, 0xFF, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x01, 0x80, 0x06, 0x00, 0x07, 0xFF, 0xFF, 0xFF, 0xFC, 0x00, 0x00,
	0x00, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x00, 0x80,
	0x04, 0x00, 0x07, 0xFF, 0xFF, 0xFF, 0xF0, 0x00, 0x38, 0x00, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x00, 0x40, 0x08, 0x00, 0x07, 0xFF, 0xFF, 0xFF, 0xE0,
	0x00, 0xFE, 0x00, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0,
	0x00, 0x40, 0x08, 0x00, 0x07, 0xFF, 0xFF, 0xFF, 0xC0, 0x03, 0xFF, 0x80, 0x00, 0xFF, 0x8F, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x00, 0x20, 0x10, 0x00, 0x07, 0xFF, 0xFF,
	0xFF, 0x80, 0x0F, 0xFF, 0xF8, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1,
	0xFF, 0xF0, 0x00, 0x30, 0x30, 0x00, 0x07, 0xFF, 0xFF, 0xFF, 0x00, 0x3F, 0xFF, 0xFC, 0x00, 0xFF,
	0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x00, 0x10, 0x20, 0x00, 0x07,
	0xFF, 0xFF, 0xFE, 0x00, 0xFF, 0xE7, 0xFE, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xF1, 0xFF, 0xF0, 0x00, 0x18, 0x40, 0x00, 0x07, 0xFF, 0xFF, 0xFC, 0x01, 0xFF, 0x81, 0xFF,
	0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x00, 0x08, 0xC0,
	0x00, 0x07, 0xFF, 0xFF, 0xFC, 0x01, 0xFF, 0x00, 0x7F, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x00, 0x04, 0x80, 0x00, 0x07, 0xFF, 0xFF, 0xF8, 0x01, 0xFF,
	0xC0, 0x7F, 0x80, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x00,
	0x05, 0x00, 0x00, 0x07, 0xFF, 0xFF, 0xF8, 0x00, 0xFF, 0xF8, 0x00, 0x00, 0xFF, 0x8F, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x00, 0x03, 0x00, 0x00, 0x07, 0xFF, 0xFF, 0xF0,
	0x00, 0x1F, 0xFE, 0x00, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF,
	0xF0, 0x00, 0x03, 0x00, 0x00, 0x07, 0xFF, 0xFF, 0xF0, 0x00, 0x07, 0xFF, 0x80, 0x00, 0xFF, 0x8F,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x00, 0x01, 0x00, 0x00, 0x07, 0xFF,
	0xFF, 0xF1, 0xFF, 0x81, 0xFF, 0x80, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xF1, 0xFF, 0xF0, 0x00, 0x01, 0x80, 0x00, 0x07, 0xFF, 0xFF, 0xE0, 0xFF, 0x00, 0x7F, 0x80, 0x00,
	0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x47, 0x80, 0x80, 0x00,
	0x07, 0xFF, 0xFF, 0xE0, 0x7F, 0xC1, 0x3F, 0x80, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x4F, 0xC0, 0x40, 0x00, 0x07, 0xFF, 0xFF, 0xE0, 0x3F, 0xF9, 0x4F,
	0x00, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x4C, 0x40,
	0x40, 0x00, 0x07, 0xFF, 0xFF, 0xE0, 0x1F, 0xFE, 0x78, 0x00, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x64, 0xC0, 0x40, 0x00, 0x07, 0xFF, 0xFF, 0xE0, 0x0F,
	0xFF, 0x80, 0x00, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0,
	0x7F, 0xC0, 0x40, 0x00, 0x07, 0xFF, 0xFF, 0xE0, 0x05, 0xFF, 0xF0, 0x00, 0x00, 0xFF, 0x8F, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x1F, 0x00, 0xC0, 0x00, 0x07, 0xFF, 0xFF,
	0xE0, 0x00, 0x3F, 0xFC, 0x00, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1,
	0xFF, 0xF0, 0x00, 0x00, 0x80, 0x00, 0x07, 0xFF, 0xFF, 0xE0, 0x00, 0x0F, 0xFF, 0x00, 0x00, 0xFF,
	0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x00, 0x01, 0x00, 0x00, 0x07,
	0xFF, 0xFF, 0xE0, 0x00, 0x03, 0xFF, 0x80, 0x00, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xF1, 0xFF, 0xF0, 0x00, 0x01, 0x00, 0x00, 0x07, 0xFF, 0xFF, 0xE0, 0x00, 0x07, 0xFF, 0x80,
	0x01, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x00, 0x02, 0x00,
	0x00, 0x07, 0xFF, 0xFF, 0xE0, 0x00, 0x3F, 0xFF, 0x00, 0x01, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x00, 0x06, 0x00, 0x00, 0x07, 0xFF, 0xFF, 0xE0, 0x00, 0xFF,
	0xFC, 0x00, 0x01, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x00,
	0x05, 0x00, 0x00, 0x07, 0xFF, 0xFF, 0xE0, 0x01, 0xFF, 0xE0, 0x00, 0x01, 0xFF, 0x8F, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x00, 0x0D, 0x80, 0x00, 0x07, 0xFF, 0xFF, 0xE0,
	0x01, 0xFF, 0x00, 0x00, 0x03, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF,
	0xF0, 0x00, 0x08, 0x80, 0x00, 0x07, 0xFF, 0xFF, 0xE0, 0x01, 0xFF, 0x80, 0x00, 0x03, 0xFF, 0x8F,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x00, 0x10, 0x40, 0x00, 0x07, 0xFF,
	0xFF, 0xE0, 0x00, 0xFF, 0xF0, 0x00, 0x03, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xF1, 0xFF, 0xF0, 0x00, 0x10, 0x40, 0x00, 0x07, 0xFF, 0xFF, 0xE0, 0x00, 0x1F, 0xFE, 0x00, 0x07,
	0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x00, 0x20, 0x20, 0x00,
	0x07, 0xFF, 0xFF, 0xE0, 0x00, 0x03, 0xFF, 0x80, 0x07, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x00, 0x60, 0x30, 0x00, 0x07, 0xFF, 0xFF, 0xE0, 0x00, 0x07, 0xFF,
	0x80, 0x0F, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x00, 0x40,
	0x10, 0x00, 0x07, 0xFF, 0xFF, 0xE0, 0x00, 0x3F, 0xFF, 0x80, 0x1F, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x00, 0xC0, 0x08, 0x00, 0x07, 0xFF, 0xFF, 0xE0, 0x01,
	0xFF, 0xFF, 0x00, 0x1F, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0,
	0x00, 0x80, 0x0C, 0x00, 0x07, 0xFF, 0xFF, 0xE0, 0x01, 0xFF, 0xF8, 0x00, 0x3F, 0xFF, 0x8F, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x01, 0x80, 0x04, 0x00, 0x07, 0xFF, 0xFF,
	0xE0, 0x01, 0xFF, 0xE0, 0x00, 0x7F, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1,
	0xFF, 0xF0, 0x01, 0x00, 0x02, 0x00, 0x07, 0xFF, 0xFF, 0xE0, 0x01, 0xFE, 0x00, 0x00, 0xFF, 0xFF,
	0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x02, 0x00, 0x03, 0x00, 0x07,
	0xFF, 0xFF, 0xE0, 0x01, 0xF0, 0x00, 0x03, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xF1, 0xFF, 0xF0, 0x06, 0x00, 0x01, 0x00, 0x07, 0xFF, 0xFF, 0xE0, 0x01, 0x80, 0x00, 0x07,
	0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x04, 0x00, 0x01,
	0x80, 0x07, 0xFF, 0xFF, 0xE0, 0x00, 0x00, 0x00, 0x1F, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x0C, 0x00, 0x00, 0x80, 0x07, 0xFF, 0xFF, 0xE0, 0x00, 0x00,
	0x00, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x08,
	0x00, 0x00, 0x40, 0x07, 0xFF, 0xFF, 0xE0, 0x00, 0x00, 0x0F, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x18, 0x00, 0x00, 0x60, 0x07, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF,
	0xF0, 0x10, 0x00, 0x00, 0x20, 0x07, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x20, 0x00, 0x00, 0x10, 0x07, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xF1, 0xFF, 0xF0, 0x60, 0x00, 0x00, 0x18, 0x07, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF0, 0x40, 0x00, 0x00, 0x08,
	0x07, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xF1, 0xFF, 0xF8, 0xC0, 0x00, 0x00, 0x04, 0x0F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xF8, 0x80, 0x00,
	0x00, 0x04, 0x1F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x02, 0x3F, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x8F, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF1,
	0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
}
