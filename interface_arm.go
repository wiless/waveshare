package ws

import (
	"log"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all"
	_ "github.com/kidoman/embd/host/rpi"
)

/// SPI related functions
var spibus embd.SPIBus

func WriteBytes(data []byte) {
	spibus.Write(data)
}

func CloseHW() {
	embd.CloseSPI()
	spibus.Close()

}

func InitHW() {
	initGPIO()
	initSPI()
}

func initGPIO() {
	embd.SetDirection(RST_PIN, embd.Out)
	embd.SetDirection(DC_PIN, embd.Out)
	embd.SetDirection(CS_PIN, embd.Out)
	embd.SetDirection(BUSY_PIN, embd.In)

	// GPIO.setup(RST_PIN, GPIO.OUT)
	// GPIO.setup(DC_PIN, GPIO.OUT)
	// GPIO.setup(CS_PIN, GPIO.OUT)
	// GPIO.setup(BUSY_PIN, GPIO.IN)

}

func initSPI() {
	//     SPI.max_speed_hz = 2000000
	//     SPI.mode = 0b00
	// # SPI device, bus = 0, device = 0
	// SPI = spidev.SpiDev(0, 0)
	if err := embd.InitSPI(); err != nil {
		log.Println("Unable to Init SPI ", err)
	}
	spibus = embd.NewSPIBus(embd.SPIMode0, channel, speed, bpw, delay)

}

func writeCmd(cmd byte) {
	embd.DigitalWrite(DC_PIN, embd.Low)
	spibus.Write([]byte{cmd})
}

func writeData(data ...byte) {
	embd.DigitalWrite(DC_PIN, embd.High)
	spibus.Write(data)

}

// ### END OF FILE ###
