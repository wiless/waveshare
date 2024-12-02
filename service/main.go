package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/mdns"
	"github.com/llgcode/draw2d"
	ws "github.com/wiless/waveshare"
)

var PORT int

func init() {
	port := os.Getenv("PORT")
	if port == "" {
		PORT = 9090
	} else {
		if num, er := strconv.Atoi(port); er != nil {
			PORT = num
		}
	}

	initEPD()
	mDNSServer()
}

var Err = func(e error) {
	if e != nil {
		fmt.Printf("Error %v", e)
	}

}

func main() {
	g := gin.New()

	g.POST("/updateepd", func(c *gin.Context) {
		var rawbytes []byte
		rawbytes, er := io.ReadAll(c.Request.Body)
		n := len(rawbytes)
		Err(er)
		if n < 5000 {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Not enough bytes", "data": n})
			return
		}
		fmt.Printf("Size is %v", n)
		img := LoadEPDBin(rawbytes)
		epd.SetFrame(img)
		epd.DisplayFrame()

	})

	g.POST("/startclock", func(c *gin.Context) {
		seconds := c.Request.FormValue("interval")
		fmt.Printf("\nFound %v", seconds)
		updateTimeBox(0)
		epd.DisplayFrame()
	})

	g.Run(fmt.Sprintf(":%d", PORT))

}

var epd ws.EPD

func initEPD() {
	draw2d.SetFontFolder(".")

	ws.InitHW()
	epd.Init(true)
}

// {
// 		bimg := Background()
// 	// ws.AsciiPrintByteImage("Background", bimg)
// 	epd.SetFrame(bimg)
// 	time.Sleep(2 * time.Second)
// 	epd.SetFrame(bimg)
// 	for {
// 		// time.Sleep(1 * time.Second)
// 		// updateTime()
// 		updateTimeBox()
// 		epd.DisplayFrame()
// 	}
// }
var server *mdns.Server

func mDNSServer() {
	// Setup our service export
	host, _ := os.Hostname()
	info := []string{"ePaper Display", "/startclock", "/updateepd"}
	service, _ := mdns.NewMDNSService(host, "_epaper._tcp", "", "", PORT, nil, info)
	log.Print("Started mDNSserver...")
	// Create the mDNS server, defer shutdown
	server, _ = mdns.NewServer(&mdns.Config{Zone: service})

}
