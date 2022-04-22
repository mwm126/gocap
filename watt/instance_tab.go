package watt

import (
	"image/color"
	"log"
	"time"

	// "aeolustec.com/capclient/cap"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type InstanceTab struct {
	TabItem     *container.TabItem
	filterEntry *widget.Entry
	table       *widget.Table
	instances   map[string][]Instance // All instances
	inst_table  [][]string            // Filtered instances
	lister      InstanceLister
	closed      bool
}

func (t *InstanceTab) Close() {
	t.closed = true
}

func NewInstanceTab(lister InstanceLister) *InstanceTab {
	t := InstanceTab{
		TabItem: nil,
		table:   nil,
		// instances:  make(map[string][]Instance),
		lister: lister,
		closed: false,
	}
	t.table = widget.NewTable(
		func() (int, int) {
			num_rows := len(t.inst_table) + 1 // add one for header
			return num_rows, 4
		},
		func() fyne.CanvasObject {
			obj := canvas.NewText("lorem ipsum", theme.PrimaryColorNamed("yellow"))
			return obj
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row == 0 {
				// give me some header
				txt := map[int]string{
					0: "Project",
					1: "UUID",
					2: "Name",
					3: "State",
				}[i.Col]
				o.(*canvas.Text).Text = txt
				o.(*canvas.Text).TextStyle.Italic = true
				o.(*canvas.Text).Color = theme.PrimaryColorNamed("gray")
				return
			}
			if i.Row-1 < len(t.inst_table) {
				o.(*canvas.Text).Color = color.White
				o.(*canvas.Text).TextStyle.Italic = false
				o.(*canvas.Text).Text = t.inst_table[i.Row-1][i.Col]
			}
		})
	t.table.SetColumnWidth(0, 200)
	t.table.SetColumnWidth(1, 500)
	t.table.SetColumnWidth(2, 200)
	t.table.SetColumnWidth(3, 200)
	t.table.Resize(fyne.NewSize(1000, 1000))

	go func() {
		for !t.closed {
			t.refresh(t.filterEntry.Text)
			time.Sleep(9 * time.Second) // TODO: configure refresh interval
		}
	}()

	scroll := container.NewScroll(t.table)
	filter_label := widget.NewLabel("Filter:")
	t.filterEntry = widget.NewEntry()
	t.filterEntry.SetPlaceHolder("<case insensitive search>")
	t.filterEntry.OnChanged = func(txt string) {
		t.inst_table = filter_instances(t.instances, txt)
		t.table.Refresh()
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
	t.inst_table = filter_instances(t.instances, txt)
	log.Printf("Refreshed: found %d instances.\n", len(t.inst_table))

	t.table.Refresh()
}
