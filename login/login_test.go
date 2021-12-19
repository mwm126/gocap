package login

import (
	"encoding/hex"
	"log"
	"net"
	"testing"

	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/cap/sshtunnel"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func NewFakeClient(server net.IP, user, pass string) (cap.Client, error) {
	client := FakeClient{}
	return &client, nil
}

type FakeClient struct{}

func (fsc FakeClient) CleanExec(command string) (string, error) {
	return "", nil
}

func (fsc FakeClient) Close() {
}

func (client FakeClient) CheckPasswordExpired(
	pass string,
	pw_expired_cb func(cap.Client),
	ch chan string,
) error {
	return nil
}

func (sc FakeClient) OpenSSHTunnel(
	user, pass string,
	local_port int,
	remote_addr string,
	remote_port int,
) sshtunnel.SSHTunnel {
	return *sshtunnel.NewSSHTunnel(
		nil,
		"testuser@localhost",
		ssh.Password(pass),
		"rem_addr:123",
		"123",
	)
}

type FakeYubikey struct {
	Available bool
}

func (yk *FakeYubikey) YubikeyAvailable() bool {
	return true
}

func (yk *FakeYubikey) FindSerial() (int32, error) {
	return 5417533, nil
}

func (yk *FakeYubikey) ChallengeResponse(chal [6]byte) ([16]byte, error) {
	if hex.EncodeToString(chal[:]) != "d459c24da2f9" {
		log.Fatal("FakeYubikey expects hardcoded challenge...", "d459c24da2f9")
	}
	var resp [16]byte
	r, _ := hex.DecodeString("9e7244e281d1e3b93f1005ba138b8a04")
	copy(resp[:], r)
	return resp, nil
}

func (yk *FakeYubikey) ChallengeResponseHMAC(chal cap.SHADigest) ([20]byte, error) {
	if hex.EncodeToString(
		chal[:],
	) != "72542b8786762da3178a035eb5f2fcef2d020dd18be729f6f67fa46ee134d5c7" {
		log.Fatal(
			"FakeYubikey expects hardcoded HMAC challenge...",
			"72542b8786762da3178a035eb5f2fcef2d020dd18be729f6f67fa46ee134d5c7",
		)
	}
	var hmac [20]byte
	h, _ := hex.DecodeString("50d30849c0e623f20665267b02fd37f4528f8cf2")
	copy(hmac[:], h)
	return hmac, nil
}

func TestLoginTab(t *testing.T) {
	a := test.NewApp()
	a.Run()
	var w fyne.Window
	connctd := container.NewVBox(widget.NewLabel("Connected!"))
	service := Service{}
	knk := cap.NewKnocker(&FakeYubikey{}, 0)
	conn_man := cap.NewCapConnectionManager(NewFakeClient, knk)
	tabs := container.NewAppTabs()
	login_tab := NewLoginTab("Login", "NETL SuperComputer", service, conn_man,
		func(login_info LoginInfo, services []Service) {
			ct := NewCapTab("test tab", "for testing", Service{},
				conn_man, func(cap *cap.Connection) {},
				connctd, login_info)
			tabs.Append(ct.Tab)
			w.SetContent(tabs)
		}, connctd, "", "")
	w = test.NewWindow(tabs)

	t.Run("Enabled", func(t *testing.T) {
		assert.False(t, login_tab.LoginForm.NetworkSelect.Disabled())
		assert.False(t, login_tab.LoginForm.UsernameEntry.Disabled())
		assert.False(t, login_tab.LoginForm.PasswordEntry.Disabled())
		assert.False(t, login_tab.LoginForm.LoginButton.Disabled())
	})

	t.Run("Disabled", func(t *testing.T) {
		login_tab.LoginForm.setEnabled(false)
		assert.True(t, login_tab.LoginForm.NetworkSelect.Disabled())
		assert.True(t, login_tab.LoginForm.UsernameEntry.Disabled())
		assert.True(t, login_tab.LoginForm.PasswordEntry.Disabled())
		assert.True(t, login_tab.LoginForm.LoginButton.Disabled())

		login_tab.LoginForm.setEnabled(true)
		assert.False(t, login_tab.LoginForm.NetworkSelect.Disabled())
		assert.False(t, login_tab.LoginForm.UsernameEntry.Disabled())
		assert.False(t, login_tab.LoginForm.PasswordEntry.Disabled())
		assert.False(t, login_tab.LoginForm.LoginButton.Disabled())
	})

	test.Tap(login_tab.LoginForm.LoginButton)

	login_tab.CloseConnection()

}
