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
	buttons            map[int]*widget.Button // for testing
}

func (t *InstanceTab) Close() {
	t.closed = true
}

func NewInstanceTab(lister InstanceLister, client SpiceClient) *InstanceTab {
	t := InstanceTab{
		TabItem:      nil,
		list:         nil,
		lister:       lister,
		spice_client: client,
		closed:       false,
		buttons:      make(map[int]*widget.Button),
	}
	t.list = widget.NewList(
		func() int {
			return len(t.filtered_instances)
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
			if i >= len(t.filtered_instances) {
				return
			}
			inst := t.filtered_instances[i]
			content := o.(*fyne.Container).Objects
			connect_btn := content[0].(*widget.Button)
			project_lbl := content[1].(*widget.Label)
			uuid_lbl := content[2].(*widget.Label)
			name_lbl := content[3].(*widget.Label)
			state_lbl := content[4].(*widget.Label)
			t.buttons[i] = connect_btn

			connect_btn.OnTapped = func() {
				port, err := t.spice_client.connect(inst)
				if err != nil {
					log.Println("Unable to connect to instance")
					return
				}
				t.spice_client.RunSpice(port)
			}

			project_lbl.SetText(inst.Project)
			uuid_lbl.SetText(inst.UUID)
			name_lbl.SetText(inst.Name)
			state_lbl.SetText(inst.State)
		})

	scroll := container.NewScroll(t.list)
	filter_label := widget.NewLabel("Filter:")
	t.filterEntry = widget.NewEntry()
	t.filterEntry.SetPlaceHolder("<case insensitive search>")
	t.filterEntry.OnChanged = t.filter
	filter := container.NewBorder(nil, nil, filter_label, nil, t.filterEntry)
	box := container.NewBorder(filter, nil, nil, nil, scroll)
	t.TabItem = container.NewTabItem("Instances", box)

	return &t
}

func (t *InstanceTab) start() {
	go func() {
		for !t.closed {
			t.refresh()
			time.Sleep(9 * time.Second) // TODO: configure refresh interval
		}
	}()
}

func (t *InstanceTab) refresh() {
	t.instances = t.lister.find_instances()
	t.filter(t.filterEntry.Text)
}

func (t *InstanceTab) filter(txt string) {
	t.filtered_instances = filter_instances(t.instances, txt)
	t.list.Refresh()
}
