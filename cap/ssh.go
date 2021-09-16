package cap

import (
	"fmt"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func newSsh(conn_man *CapConnectionManager) *container.TabItem {
	ssh := widget.NewButton("New SSH Session", func() { run_ssh(conn_man) })
	label := widget.NewLabel(
		fmt.Sprintf("or run in a terminal:  ssh localhost -p %d", SSH_LOCAL_PORT),
	)
	box := container.NewVBox(
		widget.NewLabel("To create new Terminal (SSH) Session in gnome-terminal:"),
		ssh,
		label,
	)
	return container.NewTabItem("SSH", box)
}
