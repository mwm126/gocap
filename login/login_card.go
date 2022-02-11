package login

import (
	"errors"
	"log"
	"time"

	"aeolustec.com/capclient/cap"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// A LoginCard represents the screens (Login|Connect)-(Connecting...wait or cancel)-(Password Expired...change password)-(Connected-Joule/Watt/etc.)
type LoginCard struct {
	connectionManager *cap.ConnectionManager
	Card              *widget.Card
	LoginForm         *fyne.Container
	connecting        *fyne.Container
	change_password   *fyne.Container
	connect_cancelled bool
}

//returns a pointer to a new LoginCard
func NewLoginCard(
	conn_man *cap.ConnectionManager,
	name, desc string,
	service cap.Service,
	loginForm *fyne.Container,
) *LoginCard {
	var loginCard LoginCard

	loginCard = LoginCard{
		conn_man,
		widget.NewCard(name, desc, loginForm),
		loginForm,
		func() *fyne.Container {
			connecting := widget.NewLabel("Connecting......")
			cancel := widget.NewButton("Cancel", func() {
				loginCard.connect_cancelled = true
				loginCard.Card.SetContent(loginForm)
			})
			return container.NewVBox(connecting, cancel)
		}(),
		NewChangePassword(func(new_password string) {
			conn_man.NewPasswordChannel <- new_password
			loginCard.connect_cancelled = true
			loginCard.Card.SetContent(loginForm)
		}),
		false,
	}

	return &loginCard
}

func (c *LoginCard) handle_login(
	service cap.Service,
	login_info LoginInfo,
	connected *fyne.Container,
) (*cap.Connection, error) {
	c.Card.SetContent(c.connecting)
	ext_ip, srv_ip := service.FindAddresses(login_info.Network)
	conn, err := c.connectionManager.Connect(
		login_info.Username,
		login_info.Password,
		ext_ip,
		srv_ip,
		service.CapPort,
		service.SshPort,
		func(client cap.Client) { c.Card.SetContent(c.change_password) },
		c.connectionManager.NewPasswordChannel,
	)

	if err != nil {
		log.Println("Unable to make CAP Connection: ", err)
		c.Card.SetContent(c.LoginForm)
		c.connect_cancelled = false
		return nil, errors.New("Connection failed.")
	}

	if c.connect_cancelled {
		log.Println("CAP Connection cancelled.")
		conn.Close()
		c.connect_cancelled = false
		return nil, errors.New("Connection cancelled.")
	}

	if c.connectionManager.GetPasswordExpired() {
		return nil, errors.New("Password expired; needs changed.")
	}

	time.Sleep(1 * time.Second)
	c.Card.SetContent(connected)
	return conn, nil
}
