package cap

import (
	"fmt"
	"log"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type SaveCallback func([]string)

type PortForwardTab struct {
	app      fyne.App
	save     SaveCallback
	forwards binding.StringList
}

func newPortForwardTab(app fyne.App, fwds []string, cb SaveCallback) *container.TabItem {
	t := PortForwardTab{
		app:  app,
		save: cb,
		forwards: binding.BindStringList(
			&[]string{
				"20022:localhost:22",
				"20080:localhost:80",
			},
		),
	}
	for _, fwd := range fwds {
		t.addPortForward(fwd)
	}

	add := widget.NewButton("Add", t.showNewPortForwardDialog)

	var remove *widget.Button
	var to_be_removed widget.ListItemID
	remove = widget.NewButton("Remove", func() {
		t.removeForward(to_be_removed)
		remove.Disable()
	})
	remove.Disable()

	list := widget.NewListWithData(t.forwards,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(fwd binding.DataItem, obj fyne.CanvasObject) {
			obj.(*widget.Label).Bind(fwd.(binding.String))
		})
	list.OnUnselected = func(id widget.ListItemID) { remove.Disable() }
	list.OnSelected = func(id widget.ListItemID) {
		if id < 2 {
			remove.Disable()
			log.Println("Cannot remove fixed forward #", id)
		} else {
			remove.Enable()
			to_be_removed = id
		}
	}

	box := container.NewBorder(add, remove, nil, nil, list)
	return container.NewTabItem("Port Forwards", box)
}

func (t *PortForwardTab) showNewPortForwardDialog() {
	win := t.app.NewWindow("Add Port Forward")

	local_p := widget.NewEntry()
	local_p.SetPlaceHolder("Local Port")
	remote_h := widget.NewEntry()
	remote_h.SetPlaceHolder("Remote Host")
	remote_p := widget.NewEntry()
	remote_p.SetPlaceHolder("Remote Port")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Local Port", Widget: local_p},
			{Text: "Remote Host", Widget: remote_h},
			{Text: "Remote Port", Widget: remote_p},
		},
		OnSubmit: func() {
			new_fwd := fmt.Sprintf("%s:%s:%s", local_p.Text, remote_h.Text, remote_p.Text)
			t.addPortForward(new_fwd)
			fwds, _ := t.forwards.Get()
			t.save(fwds)
			win.Close()
		},
		OnCancel:   func() { win.Close() },
		SubmitText: "Ok",
		CancelText: "Cancel",
	}
	win.SetContent(form)
	win.Show()
}

func (t *PortForwardTab) addPortForward(fwd string) {
	err := t.forwards.Append(fwd)
	if err != nil {
		log.Println("Unable to add port forward ", fwd, " because: ", err)
	}
}

func (t *PortForwardTab) removeForward(to_be_removed int) {
	fwds, _ := t.forwards.Get()
	for i := range fwds {
		if i == to_be_removed {
			fwds = append(fwds[:i], fwds[i+1:]...)
			break
		}
	}
	err := t.forwards.Set(fwds)
	if err != nil {
		log.Println("Unable to remove port forward: ", err)
	}
	t.save(fwds)
}
