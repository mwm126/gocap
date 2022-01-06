package watt

import (
	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/config"
	"aeolustec.com/capclient/forwards"
	"aeolustec.com/capclient/login"
	"aeolustec.com/capclient/ssh"
	"fmt"
	"log"

	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type WattTab struct {
	app         fyne.App
	Tabs        *container.AppTabs
	CapTab      *login.CapTab
	instanceTab *InstanceTab
}

func NewWattConnected(
	app fyne.App,
	service login.Service,
	conn_man *cap.ConnectionManager,
	login_info login.LoginInfo,
) WattTab {
	var watt_tab WattTab
	tabs := container.NewAppTabs()
	cont := container.NewMax(tabs)

	watt_tab = WattTab{
		app,
		tabs,
		login.NewCapTab("Watt", "NETL SuperComputer", service, conn_man,
			func(conn *cap.Connection) {
				watt_tab.Connect(conn)
			}, cont, login_info),
		nil,
	}
	return watt_tab
}

func (t *WattTab) Connect(conn *cap.Connection) {
	homeTab := newWattHome(func() {
		t.instanceTab.Close()
		t.CapTab.CloseConnection()
	})
	sshTab := ssh.NewSsh(conn)

	t.instanceTab = NewInstanceTab(conn)

	cfg := config.GetConfig()
	fwdTab := forwards.NewPortForwardTab(t.app, cfg.Watt_Forwards, func(fwds []string) {
		conn.UpdateForwards(fwds)
		config.SaveForwards(fwds)
	})

	t.Tabs.SetItems([]*container.TabItem{homeTab, t.instanceTab.TabItem, sshTab, fwdTab.TabItem})
}

func newWattHome(close_cb func()) *container.TabItem {
	close := widget.NewButton("Disconnect", close_cb)
	box := container.NewVBox(widget.NewLabel("Connected!"), close)
	return container.NewTabItem("Home", box)
}

func handle_ml_open_webui(conn cap.Connection) {
	redirect_httpd_port, err := cap.FreePortFinder{}.FindPort()
	if err != nil {
		log.Println("Unable to find port for Watt web interface: ", err)
		return

	}
	// url := "http://localhost:{self.local_port.web}/"
	var port uint = 123 // self.local_port.web
	cookie := get_sessionid_cookie(conn.GetUsername(), conn.GetPassword(), port)
	start_redirect_httpd(cookie, redirect_httpd_port)
	url := fmt.Sprintf("http://localhost:%d/", redirect_httpd_port)

	// LOG.exception("Unable to automatically login to Dashboard")
	// LOG.debug(f"Opening browser to:  {url}")
	open_url(url)
}

func get_sessionid_cookie(username, passwd string, port uint) string {
	// Login to dashboard from Python to get sessionid cookie

	url := fmt.Sprintf("http://localhost:%d/dashboard/auth/login/", port)
	// headers := map[string]string{"referer": url}
	// logindata := map[string]string{"username": username, "password": passwd}
	// session := requests.Session()
	log.Println("Getting csrftoken from: ", url)
	// req := session.get(url)
	// headers["x-csrftoken"] = req.cookies["csrftoken"]
	// session.post(url, logindata, headers)
	// cookie := SimpleCookie()
	// cookie["sessionid"] = session.cookies["sessionid"]
	// log.Println("Dashboard login to: %s   Returned session cookie: %s", url, session.cookies)
	return "cookie"

}

func start_redirect_httpd(cookie string, httpd_port uint) {

	url := "http://localhost:{self.local_port.web}/dashboard/project/instances/"
	html := fmt.Sprintf(`<html><head><meta http-equiv="Refresh" content="0; url=\'%s\'" /></head></html>`, url)

	// class SimpleHTTPRequestHandler(BaseHTTPRequestHandler):
	// def do_GET(self):
	// self.send_response(200)
	// self.send_header("Set-Cookie", cookie.output(header="", sep=""))
	// self.end_headers()
	fmt.Printf("Sending webpage: %s", html)
	// self.wfile.write(html.encode())

	// httpd = HTTPServer(("localhost", httpd_port), SimpleHTTPRequestHandler)
	// httpd.timeout = 10  # expect browser to open in <10 seconds
	// Thread(target=httpd.handle_request, args=(), daemon=True).start()
}

func open_url(url string) {}
