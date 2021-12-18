package login

import (
	"aeolustec.com/capclient/cap"
	"encoding/hex"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

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
	var login_info LoginInfo
	connctd := container.NewVBox(widget.NewLabel("Connected!"))
	service := Service{}
	knk := cap.NewKnocker(&FakeYubikey{}, 0)
	conn_man := cap.NewCapConnectionManager(knk)
	tabs := container.NewAppTabs()
	login_tab := NewLoginTab("Login", "NETL SuperComputer", service, conn_man,
		func(login_info LoginInfo) {
			ct := NewCapTab("test tab", "for testing", Service{},
				conn_man, func(cap *cap.Connection) {},
				connctd, login_info)
			tabs.Append(ct.Tab)
			w.SetContent(tabs)
		}, connctd, "", "")
	w = test.NewWindow(tabs)

	t.Run("Enabled", func(t *testing.T) {
		assert.False(t, login_tab.NetworkSelect.Disabled())
		assert.False(t, login_tab.UsernameEntry.Disabled())
		assert.False(t, login_tab.PasswordEntry.Disabled())
		assert.False(t, login_tab.LoginBtn.Disabled())
	})

	t.Run("Disabled", func(t *testing.T) {
		login_tab.Disable()
		assert.True(t, login_tab.NetworkSelect.Disabled())
		assert.True(t, login_tab.UsernameEntry.Disabled())
		assert.True(t, login_tab.PasswordEntry.Disabled())
		assert.True(t, login_tab.LoginBtn.Disabled())

		login_tab.Enable()
		assert.False(t, login_tab.NetworkSelect.Disabled())
		assert.False(t, login_tab.UsernameEntry.Disabled())
		assert.False(t, login_tab.PasswordEntry.Disabled())
		assert.False(t, login_tab.LoginBtn.Disabled())
	})

	login_tab.ConnectedCallback(login_info)

	login_tab.CloseConnection()

}
