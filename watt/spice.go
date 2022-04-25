package watt

import (
	"aeolustec.com/capclient/cap"
	"log"
)

type SpiceClient interface {
	connect(Instance) (uint, error)
}

type RealSpiceClient struct {
	connection cap.Connection
	uuid2port  map[string]uint
}

func (spice RealSpiceClient) connect(inst Instance) (uint, error) {
	port, err := cap.FreePortFinder{}.FindPort()
	if err != nil {
		log.Println("Unable to find free port")
		return 0, err
	}
	spice.uuid2port[inst.UUID] = port
	spice.connection.Tunnel("local_p,remote_h,remote_p")
	return port, nil
}
