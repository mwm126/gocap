package cap

import (
	"testing"
)

func TestFindPort(t *testing.T) {
	fpf := FreePortFinder{}

	port, err := fpf.FindPort()

	if err != nil {
		t.Error(err)
	}
	if !(0 < port && port < 65536) {
		t.Error(port)
	}
}
