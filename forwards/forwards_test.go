package forwards

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestPortForwardDialogDefaultForward(t *testing.T) {
	a := test.NewApp()

	fwds := []string{"000:1.2.3.4:555"}
	pft := NewPortForwardTab(a, fwds, func(fwds []string) {})

	got, _ := pft.forwards.GetValue(2)
	want := "000:1.2.3.4:555"
	if got != want {
		t.Errorf("Got %s; want %s", got, want)
	}
}

func TestPortForwardDialogRemove(t *testing.T) {
	a := test.NewApp()

	fwds := []string{"000:1.2.3.4:555"}
	pft := NewPortForwardTab(a, fwds, func(fwds []string) {})

	pft.addPortForward("777:5.5.5.5:888")

	got := pft.forwards.Length()
	want := 4
	if got != want {
		t.Errorf("Got %d; want %d", got, want)
	}

	pft.list.Select(3)
	test.Tap(pft.remove)

	got = pft.forwards.Length()
	want = 3
	if got != want {
		t.Errorf("Got %d; want %d", got, want)
	}
}

func TestPortForwardDialogCantRemoveFixed(t *testing.T) {
	a := test.NewApp()

	fwds := []string{"000:1.2.3.4:555"}
	pft := NewPortForwardTab(a, fwds, func(fwds []string) {})

	pft.addPortForward("777:5.5.5.5:888")

	got := pft.forwards.Length()
	want := 4
	if got != want {
		t.Errorf("Got %d; want %d", got, want)
	}

	pft.list.Select(0)
	test.Tap(pft.remove)

	got = pft.forwards.Length()
	want = 4
	if got != want {
		t.Errorf("Got %d; want %d", got, want)
	}
}

func TestPortForwardForm(t *testing.T) {
	a := test.NewApp()

	fwds := []string{"000:1.2.3.4:555"}
	pft := NewPortForwardTab(a, fwds, func(fwds []string) {})

	win := test.NewWindow(nil)
	pff := pft.NewPortForwardForm(win)

	pff.Form.Items[0].Widget.(*widget.Entry).SetText("111")
	pff.Form.Items[1].Widget.(*widget.Entry).SetText("6.6.6.6")
	pff.Form.Items[2].Widget.(*widget.Entry).SetText("555")

	pff.Form.OnSubmit()

	got, _ := pft.forwards.GetValue(3)
	want := "111:6.6.6.6:555"
	if got != want {
		t.Errorf("Got %s; want %s", got, want)
	}
}
