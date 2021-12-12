package forwards

import (
	"aeolustec.com/capclient/config"
	"fyne.io/fyne/v2/test"
	"testing"
)

func TestPortForwardDialog(t *testing.T) {
	a := test.NewApp()

	cfg := config.GetConfig()
	pft := NewPortForwardTab(a, cfg.Joule_Forwards, func(fwds []string) {})

	test.Tap(pft.AddButton)
}
