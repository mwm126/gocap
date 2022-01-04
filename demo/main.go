package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"

	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/config"
	"aeolustec.com/capclient/fe261"
	"aeolustec.com/capclient/joule"
	"aeolustec.com/capclient/login"
	"aeolustec.com/capclient/watt"
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/folbricht/sshtest"
	"golang.org/x/crypto/ssh"
)

func main() {
	cfg := config.GetConfig()
	cfg.Enable_joule = true
	cfg.Enable_watt = true
	cfg.Joule_Forwards = []string{}
	config.WriteConfig(cfg)

	yk := &FakeYubikey{}
	knk := cap.NewKnocker(yk, cfg.YubikeyTimeout)
	knk.StartMonitor()
	conn_man := cap.NewCapConnectionManager(cap.NewSshClient, knk)
	a := app.New()
	w := a.NewWindow("CAP Client <DEMO; Lorem Ipsum>")

	testserver, sshport := startSshServer()
	capport := sshport // doesn't matter; ignored anyway

	services := []login.Service{
		{
			Name:    "joule",
			CapPort: uint(capport),
			SshPort: sshport,
			Networks: map[string]login.Network{
				"external": {
					ClientExternalAddress: "127.0.0.1",
					CapServerAddress:      "127.0.0.1",
				},
			},
		},
		{
			Name:    "watt",
			CapPort: uint(capport),
			SshPort: sshport,
			Networks: map[string]login.Network{
				"external": {
					ClientExternalAddress: "127.0.0.1",
					CapServerAddress:      "127.0.0.1",
				},
			},
		},
	}
	err := login.InitServices(&services)
	if err != nil {
		log.Println("Could not contact Service List server:", err)
		return
	}
	client := NewClient(a, w, cfg, conn_man, sshport)
	client.Run()
	defer testserver.Close()
}

func startSshServer() (*sshtest.Server, uint) {
	test_host_key := "test-host-key"
	if err := os.RemoveAll(test_host_key); err != nil {
		panic(err)
	}

	if err := exec.Command("ssh-keygen", "-N", "", "-f", test_host_key).Run(); err != nil {
		panic(err)
	}

	hostKey := sshtest.KeyFromFile(test_host_key, "")
	server := sshtest.NewUnstartedServer()
	server.Config = &ssh.ServerConfig{NoClientAuth: true}
	server.Config.AddHostKey(hostKey)
	server.Handler = func(ch ssh.Channel, in <-chan *ssh.Request) {
		defer ch.Close()
		req, ok := <-in
		if !ok {
			return
		}
		fmt.Printf("Received '%s' request from client", req.Type)
		response := demoReplies()[string(req.Payload)]

		req.Reply(true, []byte(response))
		sshtest.SendStatus(ch, 0)
	}

	server.Start()
	_, port, err := net.SplitHostPort(server.Listener.Addr().String())
	if err != nil {
		panic(err)
	}
	portnum, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	return server, uint(portnum)
}

func demoReplies() map[string]string {
	return map[string]string{
		"env OS_PROJECT_NAME=%s openstack server list -f csv": "STUFFFFFF",
	}
}

// Client represents the Main window of CAP client
type Client struct {
	connectionManager *cap.ConnectionManager
	Tabs              *container.AppTabs
	window            fyne.Window
	app               fyne.App
	LoginTab          *login.LoginTab
}

func NewClient(
	a fyne.App,
	w fyne.Window,
	cfg config.Config,
	conn_man *cap.ConnectionManager,
	sshPort uint,
) *Client {
	var client Client

	service := login.Service{ // TODO: placeholder for real ServiceList service
		Name:    "ServiceList",
		CapPort: 62201,
		SshPort: sshPort,
		Networks: map[string]login.Network{
			"external": {
				ClientExternalAddress: "127.0.0.1",
				CapServerAddress:      "127.0.0.1",
			},
		},
	}

	connctd := container.NewVBox(widget.NewLabel("Connected!"))

	uname, pword, _ := login.GetSavedLogin()
	login_tab := login.NewLoginTab(
		"Login",
		"NETL SuperComputer",
		service,
		conn_man,
		client.setupServices,
		connctd,
		uname,
		pword,
	)

	client = Client{conn_man, nil, w, a, login_tab}
	client.setupServices(nil, make([]login.Service, 0))
	return &client
}

func (client *Client) setupServices(login_info *login.LoginInfo, services []login.Service) {
	about_tab := container.NewTabItemWithIcon(
		"About",
		theme.HomeIcon(),
		widget.NewLabel(
			"The CAP client is used for connecting to Joule, Watt, and other systems using the CAP protocol.",
		),
	)

	tabs := container.NewAppTabs(about_tab)
	tabs.SetTabLocation(container.TabLocationLeading)

	if login_info == nil {
		tabs.Append(client.LoginTab.Tab)
	} else {
		for _, service := range services {
			if service.Name == "joule" {
				joule := joule.NewJouleConnected(
					client.app,
					client.window,
					service,
					client.connectionManager,
					*login_info,
				)
				tabs.Append(joule.CapTab.Tab)
			}
			if service.Name == "watt" {
				watt := watt.NewWattConnected(
					client.app,
					service,
					client.connectionManager,
					*login_info,
				)
				tabs.Append(watt.CapTab.Tab)
			}
			if service.Name == "fe261" {
				fe261 := fe261.NewFe261Connected(
					client.app,
					service,
					client.connectionManager,
					*login_info,
				)
				tabs.Append(fe261.CapTab.Tab)
			}
		}
	}
	client.Tabs = tabs
	client.window.SetContent(client.Tabs)
}

func (client *Client) Run() {
	client.window.ShowAndRun()
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
	var resp [16]byte
	return resp, nil
}

func (yk *FakeYubikey) ChallengeResponseHMAC(chal cap.SHADigest) ([20]byte, error) {
	var hmac [20]byte
	return hmac, nil
}
