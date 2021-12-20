package ssh

import (
	"aeolustec.com/capclient/cap"

	"fmt"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewSsh(conn *cap.Connection) *container.TabItem {
	ssh := widget.NewButton("New SSH Session", func() { run_ssh(conn) })
	label := widget.NewLabel(
		fmt.Sprintf("or run in a terminal:  ssh localhost -p %d", cap.SSH_LOCAL_PORT),
	)
	box := container.NewVBox(
		widget.NewLabel("To create new Terminal (SSH) Session in gnome-terminal:"),
		ssh,
		label,
	)
	return container.NewTabItem("SSH", box)
}
