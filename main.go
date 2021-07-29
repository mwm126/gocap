package main

import (
	"aeolustec.com/capclient/cap"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Client represents the Main window of CAP client
type Client struct {
	window        fyne.Window
	usernameEntry *widget.Entry
	passwordEntry *widget.Entry
	loginBtn      *widget.Button
	knocker       cap.Knocker
}

func newClient(knocker cap.Knocker) Client {
	a := app.New()
	w := a.NewWindow("Hello")

	username := widget.NewEntry()
	username.SetPlaceHolder("Enter username...")
	password := widget.NewEntry()
	password.SetPlaceHolder("Enter password...")
	login := widget.NewButton("Login", func() {
		knocker.Knock(username.Text, password.Text)
	})
	content := container.NewVBox(username, password, login)

	w.SetContent(container.NewVBox(content))

	return Client{w, username, password, login, knocker}
}

func main() {
	knk := cap.PortKnocker{}
	client := newClient(&knk)
	client.window.ShowAndRun()
}
