package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Client struct {
	window         fyne.Window
	username_entry *widget.Entry
	password_entry *widget.Entry
	login_btn      *widget.Button
	knocker        Knocker
}

func newClient(knocker Knocker) Client {
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
	knk := PortKnocker{}
	client := newClient(&knk)
	client.window.ShowAndRun()
}
