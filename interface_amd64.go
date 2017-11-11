// Functions related to the HARDWARE configuration
package ws

import (
	"github.com/golang/glog"
)

func WriteBytes(data []byte) {
	// dummy write
}

func CloseHW() {

}

func InitHW() {
	initGPIO()
	initSPI()
}

func initGPIO() {
	glog.Infoln("Dummy : initGPIO..")
}

func initSPI() {
	glog.Infoln("Dummy : initSPI..")
}

func writeCmd(cmd byte) {
	glog.Infoln("Dummy : WriteCmd..")
}

func writeData(data ...byte) {
	glog.Infoln("Dummy : WriteData..")
}

// ### END OF FILE ###
