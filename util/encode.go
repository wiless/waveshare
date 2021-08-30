package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"strconv"
	// ws "github.com/wiless/waveshare"
)

var Err = func(e error) {
	if e != nil {
		fmt.Printf("\nError %v", e)
	}

}

var ifname string
var ofname string

var encodeCmd *flag.FlagSet
var decodeCmd *flag.FlagSet
var threshold uint
var cmd string

func init() {

	encodeCmd := flag.NewFlagSet("encode", flag.ExitOnError)

	encodeCmd.StringVar(&ifname, "image", "", "input image filename -image=abcd.png")
	encodeCmd.StringVar(&ofname, "output", "", "epd output filename -output=abcd.epd")
	encodeCmd.UintVar(&threshold, "threshold", 128, "Threshold [0,255] for 0/1 for gray scale images -threshold=128")

	decodeCmd := flag.NewFlagSet("decode ", flag.ExitOnError)
	decodeCmd.StringVar(&ifname, "epdfile", "", "Input EPD filename -epdfile=abcd.epd")

	if len(os.Args) == 1 {
		// cmd = "encode"
		// flag.Parse()
		var fn = os.Args[0]
		fmt.Println(fn + " encode <OPTIONS>")
		encodeCmd.PrintDefaults()
		fmt.Println(fn + " decode <OPTIONS>")
		decodeCmd.PrintDefaults()
		return
	}
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "encode":
			cmd = "encode"
			encodeCmd.Parse(os.Args[2:])
			if ifname == "" {
				log.Println("No input")
				return
			}
			log.Printf("Encoding %s => %s", ifname, ofname)
			if ofname == "" {
				ofname = ifname + ".epd"
			}
		case "decode":
			cmd = "decode"
			encodeCmd.Parse(os.Args[2:])
			if ifname == "" {
				log.Println("Not input filename ")
				return
			}
			log.Printf("Decoding %s ", ifname)
		default:
			cmd = "encode"
		}

	}

}

func main() {

	switch cmd {
	case "encode":
		encode(ifname, ofname)
	case "decode":
		LoadEPDBin(ifname)
	default:
		return
	}

	// fmt.Printf("%#v", epdimg)
}

func encode(fname string, ofname string) {
	f, err := os.Open(fname)
	Err(err)
	img, e := png.Decode(f)
	Err(e)
	fmt.Printf("%#v", img.Bounds().Max)

	// bitimage := toGray(img)
	b := img.Bounds()
	gimg := image.NewGray(b)
	var cg color.Gray
	mono := true
	for r := 0; r < b.Max.Y; r++ {
		fmt.Printf("\n Row %03d : ", r)
		for c := 0; c < b.Max.X; c++ {
			oldPixel := img.At(c, r)

			// gscale, _, _, _ := color.GrayModel.Convert(oldPixel).RGBA()
			cg = color.GrayModel.Convert(oldPixel).(color.Gray)
			fmt.Printf("%#v", cg)
			// convert to monochrome
			if mono {
				if cg.Y > uint8(threshold) {
					cg.Y = 255
				} else {
					cg.Y = 0
				}

			}
			gimg.SetGray(c, r, cg)

		}
	}
	// fmt.Printf("%#v", gimg)
	epdimg := ShrinkToByteImage(gimg)
	// ws.AsciiPrintByteImage("rinku", epdimg)
	SaveEPDBin(epdimg, ofname)
}

func SaveEPDBin(img image.Gray, fname string) {
	b := img.Bounds()
	R, C := b.Max.Y, b.Max.X
	fmt.Printf("\n %s = [rows x cols] = %d,%d \n", fname, R, C)
	f, er := os.Create(fname)
	defer f.Close()
	Err(er)

	for r := 0; r < R; r++ {
		fmt.Printf("\n Row %03d : ", r)
		var rowstr = ""
		var rowbytes []byte
		for c := 0; c < C; c++ {
			clr := img.GrayAt(c, r).Y
			fmt.Printf("%08b", clr)
			rowstr += fmt.Sprintf("%08b", clr)
			rowbytes = append(rowbytes, byte(clr))
		}

		f.Write(rowbytes)
		// f.WriteString("\n" + rowstr)
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

// func toGray(img image.Image) {
// 	// newimg := image.NewRGBA(image.Rect(0, 0, 200, 200))
// 	// for r := 0; r < 200; r++ {
// 	// 	for c := 0; c < 200; c++ {
// 	// 		img.Set(c, r, color.White)
// 	// 	}
// 	// }

// 	b := img.Bounds()
// 	gimg := image.NewGray(b)
// 	var cg color.Gray
// 	mono := true
// 	for r := 0; r < b.Max.Y; r++ {
// 		for c := 0; c < b.Max.X; c++ {
// 			oldPixel := img.At(c, r)

// 			// gscale, _, _, _ := color.GrayModel.Convert(oldPixel).RGBA()
// 			cg = color.GrayModel.Convert(oldPixel).(color.Gray)

// 			// convert to monochrome
// 			if mono {
// 				if cg.Y > 0 {
// 					cg.Y = 255
// 				} else {
// 					cg.Y = 0
// 				}

// 			}
// 			gimg.SetGray(c, r, cg)

// 		}
// 	}
// }

//ShrinkToByteImage assumes binary image of size R*C and genrates R*(C/8)
func ShrinkToByteImage(img *image.Gray) (byteimg image.Gray) {
	b := img.Bounds()
	R := b.Dy()
	C := b.Dx()
	CC := C / 8 // 8pixels per byte

	// if debug
	// fmt.Printf("\nImage2Byte v2 bits to Bytes %d -> %d ( RxC = %d x %d) \n ", C, CC, R, CC)

	epdimg := image.NewGray(image.Rect(0, 0, CC, R))
	var cg color.Gray
	var bitstr string
	for r := 0; r < R; r++ {
		bc := 0
		//		fmt.Printf("\n Row %03d : ", r)
		bitstr = ""
		for c := 0; c < C; c++ {
			pix := img.GrayAt(c, R-r).Y
			clr := img.At(c, r)
			u, _, _, _ := clr.RGBA()
			_ = pix
			if u > 0 { // 0 if monochrome or 128 if gray scale
				bitstr += "1"
			} else {
				bitstr += "0"
			}
			// if r < 2 {
			// 	fmt.Println(bitstr, pix, "R G B", u, v, w)
			// }
			if len(bitstr) == 8 {
				val, e := strconv.ParseUint(bitstr, 2, 8)
				if e != nil {
					log.Println(" Some error e = ", e)
				}
				// fmt.Println("Image2Byte : ", val)
				//				fmt.Print(bitstr)
				cg.Y = byte(val)

				epdimg.SetGray(bc, r, cg)
				bc++
				bitstr = ""
			}
		}
	}
	return *epdimg
}
