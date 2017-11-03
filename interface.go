// Functions related to the HARDWARE configuration
package waveshare

import (
	"log"

	"github.com/kidoman/embd"
)

// # Pin definition
var RST_PIN = 17  // GPIO_17
var DC_PIN = 25   // GPIO_25
var CS_PIN = 8    // GPIO_8
var BUSY_PIN = 24 // GPIO_24

/// SPI related functions
var spibus embd.SPIBus

const (
	channel = 0
	speed   = 2000000
	bpw     = 8
	delay   = 0
)

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

// ### END OF FILE ###
