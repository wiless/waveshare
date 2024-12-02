package ws

import (
	"fmt"
	"log"

	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/physic"
	"periph.io/x/conn/v3/spi"
	"periph.io/x/conn/v3/spi/spireg"
	"periph.io/x/host/v3"
)

type SPIbus struct {
	port spi.Port
	c    spi.Conn
}

func (s *SPIbus) Write(data []byte) (int, error) {
	// write := []byte{0x10, 0x00}
	L := len(data)
	// log.Println("Length of data %d ", len(data))
	// read := make([]byte, len(data))
	// _ = read
	err := s.c.Tx(data, nil)
	// if err := s.c.Tx(data, read); err != nil {
	// log.Fatal(err)
	// }
	if L > 0 {
		return L, err
	}
	// if L == 0 {
	return 0, nil
	// }

}
func (s *SPIbus) Close() {
	s.Close()

}
func (s *SPIbus) Init() error {
	// Open SPI port 0.0 at 1 MHz clock speed
	// port, err := spi.Open("/dev/spidev0.0")
	var err error
	s.port, err = spireg.Open("")
	if err != nil {
		fmt.Printf("Failed to open SPI port: %v\n", err)
		return err
	}
	// defer port.Close()

	// Mode:        spi.Mode0,
	// BitsPerWord: bpw,
	// MaxSpeed:    2 * physic.MegaHertz, //channel=2000000

	s.c, err = s.port.Connect(2*physic.MegaHertz, spi.Mode0, bpw)
	return err

}

func init() {
	log.Println("Done aarch64")
}

func WriteBytes(data []byte) {
	spibus.Write(data)
}

func CloseHW() {
	// embd.CloseSPI()
	spibus.Close()

}

func InitHW() {

	initGPIO()
	initSPI()
}

func DigitalWrite(pin int, level gpio.Level) error {
	pinname := fmt.Sprintf("GPIO%d", pin)
	gpin := gpioreg.ByName(pinname)
	// p := gpioreg.ByName(gpiname)
	return gpin.Out(level)

}

func initGPIO() {
	host.Init()

}

var spibus SPIbus

func initSPI() {
	//     SPI.max_speed_hz = 2000000
	//     SPI.mode = 0b00
	// # SPI device, bus = 0, device = 0
	// SPI = spidev.SpiDev(0, 0)
	// if err := embd.InitSPI(); err != nil {
	// 	log.Println("Unable to Init SPI ", err)
	// }
	// spibus = embd.NewSPIBus(embd.SPIMode0, channel, speed, bpw, delay)

	// Configure SPI settings (mode 0, 8 bits per word)

	spibus.Init()

}

func writeCmd(cmd byte) {
	// if err := embd.DigitalWrite(DC_PIN, embd.Low); err == nil {
	if err := DigitalWrite(DC_PIN, gpio.Low); err == nil {
		// spibus.Write([]byte{cmd})
		// defSPI.Write()
		if _, e := spibus.Write([]byte{cmd}); e != nil {
			log.Printf("SPI Write Error : %v", e)
		}

	} else {
		log.Println("Error writeCmd ", err)
	}

}

func writeData(data ...byte) {

	// if err := embd.DigitalWrite(DC_PIN, embd.High); err == nil {
	if err := DigitalWrite(DC_PIN, gpio.High); err == nil {
		if _, e := spibus.Write(data); e != nil {
			log.Printf("SPI Write Error : %v", e)
		}

	} else {
		log.Println("Error writeData ", err)
	}

}

// ### END OF FILE ###
