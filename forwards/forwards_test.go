package forwards

import (
	"testing"

	"aeolustec.com/capclient/config"
	"fyne.io/fyne/v2/test"
)

func TestPortForwardDialog(t *testing.T) {
	a := test.NewApp()

	cfg := config.GetConfig()
	pft := NewPortForwardTab(a, cfg.Joule_Forwards, func(fwds []string) {})

	test.Tap(pft.AddButton)
}
