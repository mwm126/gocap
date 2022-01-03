package watt

import (
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type InstanceTab struct {
	TabItem *container.TabItem
	table   *widget.Table
	data    [][]string
}

func NewInstanceTab() *InstanceTab {
	t := InstanceTab{}

	t.table = widget.NewTable(
		func() (int, int) {
			return len(t.data), len(t.data[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(t.data[i.Row][i.Col])
		})

	refresh := widget.NewButton("Refresh", func() {
		t.data = [][]string{{"Instance", "State"},
			{"instance-name", "SHUTOFF"}}

	})
	box := container.NewBorder(nil, refresh, nil, nil, t.table)
	t.TabItem = container.NewTabItem("Instances", box)
	return &t
}
