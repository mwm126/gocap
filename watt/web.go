package watt

import (
	"fmt"
	"log"
	"net/http"

	"aeolustec.com/capclient/cap"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gorilla/sessions"
)

type WebTab struct {
	TabItem      *container.TabItem
	filterEntry  *widget.Entry
	table        *widget.Table
	instances    map[string][]Web
	inst_table   [][]string
	connection   *cap.Connection
	tunnel_spice *cap.Tunnel
	tunnel_web   *cap.Tunnel
	closed       bool
}

type Web struct {
	UUID  string
	Name  string
	State string
}

func (t *WebTab) Close() {
	t.closed = true
}

func NewWebTab(conn *cap.Connection) *WebTab {
	t := WebTab{
		TabItem:    nil,
		table:      nil,
		instances:  make(map[string][]Web),
		connection: conn,
		closed:     false,
	}
	label := widget.NewLabel("Start Web Interface:")
	btn := widget.NewButton("start web interface", t.handle_ml_open_webui)
	box := container.NewBorder(nil, nil, label, btn)
	t.TabItem = container.NewTabItem("Webs", box)

	t.tunnel_spice = fwd_spice(conn)
	t.tunnel_web = fwd_web(conn)
	return &t
}

func fwd_spice(conn *cap.Connection) *cap.Tunnel {
	tunnel, err := conn.NewTunnel(WATT_SPICE_PORT, WATT_SPICE_HOST, WATT_SPICE_PORT)
	if err != nil {
		log.Printf("Could not forward %d:%s:%d because %s", WATT_SPICE_PORT, WATT_SPICE_HOST, WATT_SPICE_PORT, err)
		return nil
	}
	return tunnel
}

func fwd_web(conn *cap.Connection) *cap.Tunnel {
	local_web_port, err := cap.FreePortFinder{}.FindPort()
	if err != nil {
		log.Println("Could not find web SPICE port: ", err)
		return nil
	}

	tunnel, err := conn.NewTunnel(local_web_port, WATT_WEB_HOST, WATT_WEB_PORT)
	if err != nil {
		log.Printf("Could not forward %d:%s:%d because %s", local_web_port, WATT_WEB_HOST, WATT_WEB_PORT, err)
		return nil
	}
	return tunnel
}

const WATT_WEB_HOST string = "192.168.101.200"
const WATT_WEB_PORT uint = 80
const WATT_SPICE_HOST string = "192.168.101.200"
const WATT_SPICE_PORT uint = 6082

func (t *WebTab) handle_ml_open_webui() {
	if t.tunnel_web == nil {
		log.Println("No Watt web tunnel available")
		return
	}
	local_port_web := t.tunnel_web.LocalPort()

	redirect_httpd_port, err := cap.FreePortFinder{}.FindPort()
	if err != nil {
		log.Println("Unable to find port for Watt web interface: ", err)
		return
	}

	cookie := get_sessionid_cookie(t.connection.GetUsername(), t.connection.GetPassword(), local_port_web)

	go start_redirect_httpd(cookie, redirect_httpd_port, local_port_web)

	url := fmt.Sprintf("http://localhost:%d/", redirect_httpd_port)
	go open_url(url)
}

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func get_sessionid_cookie(username, passwd string, port uint) *http.Cookie {
	url := fmt.Sprintf("http://localhost:%d/dashboard/auth/login/", port)
	headers := map[string]string{"referer": url}
	// logindata := map[string]string{"username": username, "password": passwd}
	sess := sessions.NewSession(store, "watt-session")
	log.Println("Getting csrftoken from: ", url)
	req := sess.Options.Path
	headers["x-csrftoken"] = req
	// sess.post(url, logindata, headers)
	cookie := sessions.NewCookie("cookiename", "cookievalue", store.Options)
	// cookie["sessionid"] = sess.cookies["sessionid"]
	return cookie
}

func start_redirect_httpd(cookie *http.Cookie, httpd_port uint, local_port_web uint) {
	url := fmt.Sprintf("http://localhost:%d/dashboard/project/instances/", local_port_web)
	html := fmt.Sprintf(`<html><head><meta http-equiv="Refresh" content="0; url=\'%s\'" /></head></html>`, url)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Setting Cookie: ", *cookie)
		http.SetCookie(w, cookie)
		w.WriteHeader(200)
		log.Println("Sending HTML: ", html)
		fmt.Fprintf(w, html)
	})
	log.Println("Going to listen and serve on port:  ", httpd_port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", httpd_port), nil); err != nil {
		log.Printf("could not listen on Spice httpd port %d: %s", httpd_port, err)
	}
}
