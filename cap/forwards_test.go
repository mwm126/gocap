package cap

import (
	"fyne.io/fyne/v2/test"
	"testing"
)

func TestPortForwardDialog(t *testing.T) {
	a := test.NewApp()

	cfg := GetConfig()
	pft := NewPortForwardTab(a, cfg.Joule_Forwards, func(fwds []string) {})

	test.Tap(pft.AddButton)
}
