package watt

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"aeolustec.com/capclient/cap"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type WebTab struct {
	TabItem      *container.TabItem
	table        *widget.Table
	instances    map[string][]Web
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

	// url := fmt.Sprintf("http://localhost:%d/", redirect_httpd_port)
	url := fmt.Sprintf("http://localhost:%d/dashboard/project/instances/", local_port_web)
	go open_url(url)
}

func start_redirect_httpd(cookie *http.Cookie, httpd_port uint, local_port_web uint) {
	url := fmt.Sprintf("http://localhost:%d/dashboard/project/instances/", local_port_web)
	html := fmt.Sprintf(`<html><head><meta http-equiv="Refresh" content="0; url=\'%s\'" /></head></html>`, url)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Setting Cookie: ", *cookie)
		http.SetCookie(w, cookie)
		w.WriteHeader(200)
		log.Println("Sending HTML: ", html)
		fmt.Fprint(w, html)
	})
	log.Println("Going to listen and serve on port:  ", httpd_port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", httpd_port), nil); err != nil {
		log.Printf("could not listen on Spice httpd port %d: %s", httpd_port, err)
	}
}

func get_sessionid_cookie(username, passwd string, port uint) *http.Cookie {
	req_url := fmt.Sprintf("http://localhost:%d/dashboard/auth/login/", port)
	logindata := url.Values{"username": {username}, "password": {passwd}}
	log.Println("Getting csrftoken from: ", req_url)

	resp, err := http.Get(req_url)
	log.Println("GET header:   ", resp.Header)
	if err != nil {
		panic(err)
	}
	log.Println("GET Cookies:", resp.Cookies())
	if len(resp.Cookies()) != 1 {
		log.Println("Expected one cookie in response: ", req_url)
		return nil
	}
	if resp.Cookies()[0].Name != "csrftoken" {
		log.Println("Expected csrftoken cookie in response: ", req_url)
		return nil
	}
	csrftoken := resp.Cookies()[0].Value

	req, err := http.NewRequest("POST", req_url, strings.NewReader(logindata.Encode()))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	req.Header.Set("referer", req_url)
	req.Header.Set("x-csrftoken", csrftoken)

	// Create and Add cookie to request
	cookie := http.Cookie{Name: "cookie_name", Value: "cookie_value"}
	req.AddCookie(&cookie)

	// Set client timeout
	client := &http.Client{Timeout: time.Second * 10}

	// Validate cookie and headers are attached
	fmt.Println(req.Cookies())
	fmt.Println(req.Header)

	// Send request
	resp, err = client.Do(req)
	if err != nil {
		log.Println("Error with request: ", err)
	}
	log.Println("Body: ", resp.Body)

	// Read body from response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	fmt.Printf("%s\n", body)
	defer resp.Body.Close()
	var sessionid string

	for _, cookie := range resp.Cookies() {
		log.Println("cookie: ", cookie.Name, "     ---     ", cookie.Value)
		if cookie.Name == "sessionid" {
			sessionid = cookie.Value
		}
	}
	fmt.Println("sessionid:", sessionid)
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)

	cooookie := resp.Cookies()[0]
	log.Printf("Dashboard login to: %s   Returned session cookie: %s", req_url, cooookie)
	log.Println("===================COOKIE+++++++++++++++++++")
	return cooookie
}

// here is your cookie [header], pinhead.  sessionid=cta2otqcb6va50xvlbspny9z64977t0u
// 127.0.0.1 - - [09/Jan/2022 08:39:02] "GET / HTTP/1.1" 200 -
//	Sending HTML:  b'<html><head><meta http-equiv="Refresh" content="0; url=\'http://localhost:50789/dashboard/project/instances/\'" /></head></html>'
