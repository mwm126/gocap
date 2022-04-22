package watt

import (
	"log"
	// "time"
	// // "aeolustec.com/capclient/cap"
	// fyne "fyne.io/fyne/v2"
	// "fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/layout"
	// "fyne.io/fyne/v2/widget"
)

type SpiceClient interface {
	connect(Instance)
}

type RealSpiceClient struct{}

func (spice RealSpiceClient) connect(inst Instance) {
	log.Println("CONNECT TO INSTANCE ", inst.UUID)
}
