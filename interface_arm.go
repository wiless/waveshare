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
	log.Panicln("ARM mode..")
}

func WriteBytes(data []byte) {
	spibus.Write(data)
}

func CloseHW() {
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

	host.Init() // moving periph.io ( suppoerts recent kernel)

}

var spibus SPIbus

func initSPI() {
	spibus.Init()

}

func writeCmd(cmd byte) {
	DigitalWrite(DC_PIN, gpio.Low)
	spibus.Write([]byte{cmd})
}

func writeData(data ...byte) {
	DigitalWrite(DC_PIN, gpio.High)
	spibus.Write(data)

}

// ### END OF FILE ###
