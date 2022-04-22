package watt

import (
	"log"
	"time"

	// "aeolustec.com/capclient/cap"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type InstanceTab struct {
	TabItem            *container.TabItem
	filterEntry        *widget.Entry
	list               *widget.List
	instances          map[string][]Instance // All instances
	filtered_instances []Instance            // Filtered instances
	lister             InstanceLister
	spice_client       SpiceClient
	closed             bool
}

func (t *InstanceTab) Close() {
	t.closed = true
}

func NewInstanceTab(lister InstanceLister) *InstanceTab {
	t := InstanceTab{
		TabItem: nil,
		list:    nil,
		// instances:  make(map[string][]Instance),
		lister: lister,
		closed: false,
	}
	t.list = widget.NewList(
		func() int {
			num_rows := len(t.filtered_instances)
			return num_rows
		},
		func() fyne.CanvasObject {
			connect_btn := widget.NewButton("Connect", func() {})
			project_lbl := widget.NewLabel("")
			uuid_lbl := widget.NewLabel("")
			name_lbl := widget.NewLabel("")
			state_lbl := widget.NewLabel("")

			content := container.New(layout.NewHBoxLayout(), connect_btn, project_lbl, uuid_lbl, name_lbl, state_lbl)

			return content
		},
		func(i int, o fyne.CanvasObject) {
			inst := t.filtered_instances[i]
			content := o.(*fyne.Container).Objects
			connect_btn := content[0].(*widget.Button)
			project_lbl := content[1].(*widget.Label)
			uuid_lbl := content[2].(*widget.Label)
			name_lbl := content[3].(*widget.Label)
			state_lbl := content[4].(*widget.Label)

			connect_btn.OnTapped = func() {
				t.spice_client.connect(inst)
				RunSpice(12345)
			}
			project_lbl.SetText(inst.Project)
			uuid_lbl.SetText(inst.UUID)
			name_lbl.SetText(inst.Name)
			state_lbl.SetText(inst.State)
		})

	go func() {
		for !t.closed {
			t.refresh(t.filterEntry.Text)
			time.Sleep(9 * time.Second) // TODO: configure refresh interval
		}
	}()

	scroll := container.NewScroll(t.list)
	filter_label := widget.NewLabel("Filter:")
	t.filterEntry = widget.NewEntry()
	t.filterEntry.SetPlaceHolder("<case insensitive search>")
	t.filterEntry.OnChanged = func(txt string) {
		log.Println("t.instances:::: ", len(t.instances))
		t.filtered_instances = filter_instances(t.instances, txt)
		log.Println("t.filtered_instances:::: ", len(t.filtered_instances))
		t.list.Refresh()
	}
	filter := container.NewBorder(nil, nil, filter_label, nil, t.filterEntry)
	box := container.NewBorder(filter, nil, nil, nil, scroll)
	t.TabItem = container.NewTabItem("Instances", box)
	return &t
}

func (t *InstanceTab) refresh(txt string) {
	t.instances = t.lister.find_instances()
	// if err != nil {
	// log.Println("Could not refresh", err)
	// return
	// } else {
	// t.instances = instmap
	// }
	t.filtered_instances = filter_instances(t.instances, txt)
	log.Printf("Refreshed: found %d instances.\n", len(t.instances))

	t.list.Refresh()
}
