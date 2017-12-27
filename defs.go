// Package ws implements functions  to communication, commands protocols for the EPD
package ws

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/golang/glog"
	"golang.org/x/image/font/gofont/goregular"

	"github.com/kidoman/embd"
)

// # Display resolution
var EPD_WIDTH uint8 = 200
var EPD_HEIGHT uint8 = 200

var EPD_FONT *truetype.Font

// # EPD1IN54 commands
var DRIVER_OUTPUT_CONTROL byte = 0x01
var BOOSTER_SOFT_START_CONTROL byte = 0x0C
var GATE_SCAN_START_POSITION byte = 0x0F
var DEEP_SLEEP_MODE byte = 0x10
var DATA_ENTRY_MODE_SETTING byte = 0x11
var SW_RESET byte = 0x12
var TEMPERATURE_SENSOR_CONTROL byte = 0x1A
var MASTER_ACTIVATION byte = 0x20
var DISPLAY_UPDATE_CONTROL_1 byte = 0x21
var DISPLAY_UPDATE_CONTROL_2 byte = 0x22
var WRITE_RAM byte = 0x24
var WRITE_VCOM_REGISTER byte = 0x2C
var WRITE_LUT_REGISTER byte = 0x32
var SET_DUMMY_LINE_PERIOD byte = 0x3A
var SET_GATE_TIME byte = 0x3B
var BORDER_WAVEFORM_CONTROL byte = 0x3C
var SET_RAM_X_ADDRESS_START_END_POSITION byte = 0x44
var SET_RAM_Y_ADDRESS_START_END_POSITION byte = 0x45
var SET_RAM_X_ADDRESS_COUNTER byte = 0x4E
var SET_RAM_Y_ADDRESS_COUNTER byte = 0x4F
var TERMINATE_FRAME_READ_WRITE byte = 0xFF

func init() {
	flag.Parse()
	EPD_FONT, _ = truetype.Parse(goregular.TTF)
}

type EPD struct {
	lutFull bool
	// Sequence for updating
	lutFullUpdate    []byte
	lutPartialUpdate []byte
	screen           int
}

func (e *EPD) SetDefaults() {
	e.lutFullUpdate = []byte{
		0x02, 0x02, 0x01, 0x11, 0x12, 0x12, 0x22, 0x22,
		0x66, 0x69, 0x69, 0x59, 0x58, 0x99, 0x99, 0x88,
		0x00, 0x00, 0x00, 0x00, 0xF8, 0xB4, 0x13, 0x51,
		0x35, 0x51, 0x51, 0x19, 0x01, 0x00}

	e.lutPartialUpdate = []byte{
		0x10, 0x18, 0x18, 0x08, 0x18, 0x18, 0x08, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x13, 0x14, 0x44, 0x12,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

	e.lutFull = true
}

func (e *EPD) SendCommand(cmd byte) {

	writeCmd(cmd)

}

func (e *EPD) SendData(data ...byte) {
	writeData(data...)
}

func (e *EPD) CallFunction(command byte, data ...byte) {
	e.SendCommand(command)
	e.SendData(data...)
}

func (e *EPD) Init(full bool) {
	if len(e.lutFullUpdate) == 0 || len(e.lutPartialUpdate) == 0 {
		e.SetDefaults()
	}

	var dataseq []byte
	e.lutFull = full
	// self.lut = lut
	// self.reset()
	e.reset()
	// self.send_command(DRIVER_OUTPUT_CONTROL)
	// self.send_data((EPD_HEIGHT - 1) & 0xFF)
	// self.send_data(((EPD_HEIGHT - 1) >> 8) & 0xFF)
	// self.send_data(0x00)                     # GD = 0 SM = 0 TB = 0
	dataseq = []byte{(EPD_HEIGHT - 1) & 0xFF, ((EPD_HEIGHT - 1) >> 8) & 0xFF, 0x00}
	e.CallFunction(DRIVER_OUTPUT_CONTROL, dataseq...)

	// self.send_command(BOOSTER_SOFT_START_CONTROL)
	// self.send_data(0xD7)
	// self.send_data(0xD6)
	// self.send_data(0x9D)
	e.CallFunction(BOOSTER_SOFT_START_CONTROL, 0xD7, 0xD6, 0x9D)

	// self.send_command(WRITE_VCOM_REGISTER)
	// self.send_data(0xA8)                     # VCOM 7C
	e.CallFunction(WRITE_VCOM_REGISTER, 0xA8)

	// self.send_command(SET_DUMMY_LINE_PERIOD)
	// self.send_data(0x1A)                     # 4 dummy lines per gate
	e.CallFunction(SET_DUMMY_LINE_PERIOD, 0x1A)

	// self.send_command(SET_GATE_TIME)
	// self.send_data(0x08)                     # 2us per line
	e.CallFunction(SET_GATE_TIME, 0x08)

	// self.send_command(DATA_ENTRY_MODE_SETTING)
	// self.send_data(0x03)                     # X increment Y increment
	e.CallFunction(DATA_ENTRY_MODE_SETTING, 0x03)

	e.setLookupTable(e.lutFull)

}

//reset - module reset.often used to awaken the module in deep sleep,
func (e *EPD) reset() {
	embd.DigitalWrite(RST_PIN, embd.Low)
	time.Sleep(200 * time.Millisecond)
	embd.DigitalWrite(RST_PIN, embd.High)
	time.Sleep(200 * time.Millisecond)

}

//
//   @brief: set the look-up table register
func (e *EPD) setLookupTable(full bool) {
	e.lutFull = full

	if e.lutFull {
		e.CallFunction(WRITE_LUT_REGISTER, e.lutFullUpdate...)
	} else {
		e.CallFunction(WRITE_LUT_REGISTER, e.lutPartialUpdate...)
	}

}

// Ensure to wait before any next command is executed.. monitors the
// BUSY_PIN
func (e *EPD) wait() {
	var busy int
	var err error
	for ; busy == 1; busy, err = embd.DigitalRead(BUSY_PIN) {
		if err != nil {
			log.Panic("Error waiting BUSY_PIN", err)
		}
		time.Sleep(100 * time.Millisecond) // polling for every 100ms
	}

}

// wait_until_idle(self):
//         while(self.digital_read(self.busy_pin) == 1):      # 0: idle, 1: busy
//             self.delay_ms(100)

func (e *EPD) Sleep(full bool) {
	e.CallFunction(DEEP_SLEEP_MODE)
	e.wait()
	//  self.send_command(DEEP_SLEEP_MODE)
}

func (e *EPD) Screen() int {
	return e.screen
}

// ##
//  #  @brief: update the display
//  #          there are 2 memory areas embedded in the e-paper display
//  #          but once this function is called,
//  #          the the next action of SetFrameMemory or ClearFrame will
//  #          set the other memory area.
//  ##
func (e *EPD) DisplayFrame() {
	if e.screen == 0 {
		e.screen = 1 // next frame where image will be set
	} else {
		e.screen = 0
	}
	log.Println("Current SCREEN  ", e.screen)
	e.CallFunction(DISPLAY_UPDATE_CONTROL_2, 0xC4)
	e.CallFunction(MASTER_ACTIVATION)
	e.CallFunction(TERMINATE_FRAME_READ_WRITE)
	e.wait()
}

// ##
//  #  @brief: specify the memory area for data R/W
// def set_memory_area(self, x_start, y_start, x_end, y_end)
func (e *EPD) setMemArea(x0, y0, x1, y1 byte) {
	//   x point must be the multiple of 8 or the last 3 bits will be ignored
	e.CallFunction(SET_RAM_X_ADDRESS_START_END_POSITION, (x0>>3)&0xFF, (x1>>3)&0xFF)
	e.CallFunction(SET_RAM_Y_ADDRESS_START_END_POSITION, y0&0xFF, (y0>>8)&0xFF, y1&0xFF, (y1>>8)&0xFF)
}

/*
   @brief: specify the start point for data R/W in the memory
   //set_memory_pointer()
*/
func (e *EPD) SetXY(x, y byte) {
	e.CallFunction(SET_RAM_X_ADDRESS_COUNTER, (x>>3)&0xFF)
	e.CallFunction(SET_RAM_Y_ADDRESS_COUNTER, y&0xFF, (y>>8)&0xFF)
	e.wait()
}

// #
//  #  @brief: clear the frame memory with the specified color.
//  #          this won't update the display.
func (e *EPD) ClearFrame(color byte) {
	e.setMemArea(0, 0, EPD_WIDTH-1, EPD_HEIGHT-1)
	//	e.setMemArea(0, 0, 200,200)
	e.SetXY(0, 0)
	e.SendCommand(WRITE_RAM)

	L := int((EPD_WIDTH / 8)) * int(EPD_HEIGHT) // 8pixels cols = 1 byte

	for i := 0; i < L; i++ {
		e.SendData(color)
	}
}

// ##
//  #  @brief: convert an image to a buffer
//  ## Generates a Byte Buffer
// def get_frame_buffer(self, image):
func (e *EPD) GetFrame() *image.Gray {
	img := image.NewGray(image.Rect(0, 0, int(EPD_WIDTH), int(EPD_HEIGHT)))

	return img
}

var mode bool = true
var rval uint8

func init() {
	rval = uint8(rand.Int31n(255))
}

func AsciiPrintBytes(name string, img image.Gray) {
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

// SetSubFrame sets subset of image at r,c location, assume r,c=8n , column is multiple of 8
func (e *EPD) SetSubFrame(r, c int, binimg *image.Gray) {

	W, H := binimg.Bounds().Dx(), binimg.Bounds().Dy()

	byteimg := Mono2ByteImage(binimg)
	// AsciiPrintBytes("SUBIMAGE", byteimg)
	_ = W
	BW := byteimg.Bounds().Dx()
	hh := H
	//	BW := 6 // 6*8=48 PIXEL wide

	e.setMemArea(uint8(c), uint8(r), uint8(c+BW*8-1), uint8(r+hh-1))
	//	log.Println("Rand val ", rval, W, BW)

	e.SetXY(byte(c), byte(r))

	e.SendCommand(WRITE_RAM)
	for row := 0; row < hh; row++ {
		bytearray := make([]byte, BW)

		for col := 0; col < BW; col++ {
			pixel := byteimg.GrayAt(col, row).Y
			//			pixel := 0X80
			//	pixel = 0xAA
			//pixel := rval
			//	pixel= uint8(rand.Int31n(255))
			//	if row%2 == 0 {
			//		pixel = 0xFF
			//	}
			bytearray[col] = pixel // byte(rval)
		}

		e.SendData(bytearray...)
	}

	e.wait()
	e.DisplayFrame()
}

// SetSubFrame sets subset of image at r,c location, assume r,c=8n , column is multiple of 8
func (e *EPD) FillSubFrame(r, c int, binimg *image.Gray) {

	W, H := binimg.Bounds().Dx(), binimg.Bounds().Dy()

	byteimg := Mono2ByteImage(binimg)
	// AsciiPrintBytes("SUBIMAGE", byteimg)
	_ = W
	BW := byteimg.Bounds().Dx()
	hh := H
	//	BW := 6 // 6*8=48 PIXEL wide

	e.setMemArea(uint8(c), uint8(r), uint8(c+BW*8-1), uint8(r+hh-1))
	//	log.Println("Rand val ", rval, W, BW)

	e.SetXY(byte(c), byte(r))

	e.SendCommand(WRITE_RAM)
	for row := 0; row < hh; row++ {
		bytearray := make([]byte, BW)

		for col := 0; col < BW; col++ {
			pixel := byteimg.GrayAt(col, row).Y
			//			pixel := 0X80
			//	pixel = 0xAA
			//pixel := rval
			//	pixel= uint8(rand.Int31n(255))
			//	if row%2 == 0 {
			//		pixel = 0xFF
			//	}
			bytearray[col] = pixel // byte(rval)
		}

		e.SendData(bytearray...)
	}

	e.wait()
	// e.DisplayFrame()
}

func (e *EPD) DrawLine(row int, thick int, color uint8) {

	e.setMemArea(0, byte(row), 200, byte(row+thick-1))
	bytearray := make([]byte, 25)
	e.SetXY(0, byte(row))
	for c := 0; c < 25; c++ {
		if color > 0 {
			bytearray[c] = 0xff
		}
	}
	for r := 0; r < thick; r++ {
		e.CallFunction(WRITE_RAM, bytearray...)
		e.wait()
	}

}

//  #  @brief: put an (SUB) image to the frame memory.
//  #          this won't update the display.
func (e *EPD) SetFrame(byteimg image.Gray) {
	w, h := byte(byteimg.Bounds().Dx()), byte(byteimg.Bounds().Dy())
	if h < 200 || w < 25 {
		glog.Errorln("Image large size ", h, w)
		return
	}
	// var x1, y1 byte
	// x1 = x0 + (w) - 1
	// y1 = y0 + (h) - 1
	// if x0+w >= EPD_WIDTH {
	// 	x1 = EPD_WIDTH - 1
	// }
	// if y0+h >= EPD_HEIGHT {
	// 	y1 = EPD_HEIGHT - 1
	// }

	e.setMemArea(0, 0, 200, 200)

	// # send the image data

	// rr := int(y1 - y0 + 1)
	// cc := int(x1 - x0 + 1)

	e.SetXY(0, byte(0))

	//	e.SendCommand(WRITE_RAM)
	for row := 0; row < 200; row++ {
		bytearray := make([]byte, 25)
		for col := 0; col < 25; col++ {
			pixel := byteimg.GrayAt(col, row).Y
			//			pixel := 0X80
			bytearray[col] = byte(pixel)
		}

		//	e.SendCommand(WRITE_RAM)
		//	e.SendData(bytearray...)
		e.CallFunction(WRITE_RAM, bytearray...)
		e.wait()
	}
	e.DisplayFrame()
}

func (e *EPD) WriteBytePixel(row, col byte, pixel ...byte) {
	e.SetXY(col, row)
	e.SendCommand(WRITE_RAM)
	e.SendData(pixel...)
	e.wait()

}

//Image2Byte assumes binary image of size R*C = R*(C/8)
func Mono2ByteImage(img *image.Gray) (byteimg image.Gray) {
	return Mono2ByteImagev2(img)
	R := img.Rect.Dy()
	C := img.Rect.Dx()
	CC := C / 8 // 8pixels per byte

	// if debug
	// fmt.Println("Image2Byte bits to Bytes ", C, CC)

	epdimg := image.NewGray(image.Rect(0, 0, R, CC))
	var cg color.Gray
	var bitstr string
	for r := 0; r < R; r++ {
		bc := 0
		//		fmt.Printf("\n Row %d : ", r)
		bitstr = ""
		for c := 0; c < C; c++ {
			pix := img.GrayAt(R-r, c).Y
			if pix > 0 { // 0 if monochrome or 128 if gray scale
				bitstr += "1"
			} else {
				bitstr += "0"
			}

			if len(bitstr) == 8 {
				val, e := strconv.ParseUint(bitstr, 2, 8)
				if e != nil {
					log.Println(" Some error e = ", e)
				}
				// fmt.Println("Image2Byte : ", val)
				//				fmt.Print(bitstr)
				cg.Y = byte(val)

				epdimg.SetGray(r, bc, cg)
				bc++
				bitstr = ""
			}
		}
	}
	return *epdimg
}

func logme(info string, e error) {
	if e != nil {
		log.Panicln(info, " : ", e)
	}
}

//Image2Byte assumes binary image of size R*C = R*(C/8)
func Mono2ByteImagev2(img *image.Gray) (byteimg image.Gray) {
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
