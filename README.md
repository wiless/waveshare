# waveshare
A golang library for eInk Paper display from waveshare ( www.waveshare.com)

There are lot of [ePaper displays](https://www.waveshare.com/product/modules/oleds-lcds/e-paper.htm) available online.

# Pin Diagram
```
VCC	3.3V
GND	Ground
DIN	SPI MOSI pin
CLK	SPI SCK pin
CS	SPI chip selection, low active
DC	Data/Command selection (high for data, low for command)
RST	External reset, low active
BUSY	Busy status output, high active


# PI3 
RST_PIN = 17  // GPIO_17
DC_PIN = 25   // GPIO_25
CS_PIN = 8    // GPIO_8
BUSY_PIN = 24 // GPIO_24
```
