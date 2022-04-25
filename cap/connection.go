package cap

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"aeolustec.com/capclient/cap/sshtunnel"
	"fyne.io/fyne/v2/data/binding"
)

type Client interface {
	// NewSession() (Client, error)
	CleanExec(cmd string) (string, error)
	Close()
	// CombinedOutput(cmd string) ([]byte, error)

	OpenSSHTunnel(
		user, pass string,
		local_port int,
		remote_addr string,
		remote_port int,
	) sshtunnel.SSHTunnel
	CheckPasswordExpired(string, func(Client), chan string) error
	Dial(protocol, endpoint string) (net.Conn, error)
}

// A Connection represents a successful SSH connection after the port knock
type Connection struct {
	client       Client
	forwards     map[string]sshtunnel.SSHTunnel // Set by user
	tunnels      map[string]sshtunnel.SSHTunnel // Added for SPICE
	username     string
	password     string
	uid          string
	hostName     string
	loginAddr    string
	webLocalPort int
	sshLocalPort int
}

func (c *Connection) GetUsername() string {
	return c.username
}

func (c *Connection) GetPassword() string {
	return c.password
}

func (c *Connection) GetAddress() string {
	return c.loginAddr
}

func (c *Connection) GetUid() string {
	return c.uid
}

func (c *Connection) GetHostname() string {
	return c.hostName
}

func (c *Connection) GetClient() Client {
	return c.client
}

func (conn *Connection) UpdateForwards(fwds []string) {
	for _, fwd := range fwds {
		if _, missing := conn.forwards[fwd]; missing {
			conn.forwards[fwd] = conn.forward(fwd)
		}
	}
	for forward, tunnel := range conn.forwards {
		found := false
		for _, fwd := range fwds {
			if forward == fwd {
				found = true
				break
			}
		}
		if !found {
			tunnel.Close()
			delete(conn.forwards, forward)
		}
	}
}

func (conn *Connection) Tunnel(fwd string) {
	if _, missing := conn.forwards[fwd]; missing {
		conn.tunnels[fwd] = conn.forward(fwd)
		log.Println("Tunneling: ", fwd)
		return
	}
	log.Println("Warning:  already tunneling: ", fwd)
}

func (conn *Connection) forward(fwd string) sshtunnel.SSHTunnel {

	result := strings.Split(fwd, ",")
	local_p, err := strconv.Atoi(result[0])
	if err != nil {
		log.Println("Warning:  bad local port", err)
	}
	remote_h := result[1]
	remote_p, err := strconv.Atoi(result[2])
	if err != nil {
		log.Println("Warning:  bad remote port", err)
	}
	return conn.client.OpenSSHTunnel(conn.username,
		conn.password, local_p, remote_h, remote_p)
}

func (c Connection) Close() {
	if c.client == nil {
		log.Println("Not connected; Cannot close connection")
		return
	}
	for _, fwd := range c.forwards {
		fwd.Close()
	}
	for _, tun := range c.tunnels {
		tun.Close()
	}
	c.client.Close()
}

func (c *Connection) FindSessions() ([]Session, error) {
	sessions := make([]Session, 0, 10)
	if c.client == nil {
		return sessions, errors.New("Client missing")
	}
	text, err := c.client.CleanExec("ps auxnww|grep Xvnc|grep -v grep")
	if err != nil {
		return sessions, errors.New("Unable to find sessions")
	}
	return parseSessions(c.GetUsername(), text), nil
}

const MAX_GUI_COUNT = 4

func (c *Connection) CreateVncSession(xres string, yres string) (string, string, error) {
	sessions, err := c.FindSessions()
	if err != nil {
		return "", "", err
	}
	gui_count := len(sessions)
	if MAX_GUI_COUNT <= gui_count {
		return "", "", errors.New("Users may only have {MAX_GUI_COUNT} VNC sessions open at once.")
	}

	otp, displayNumber, err := c.startVncSession(xres, yres)
	if err != nil {
		return "", "", err
	}

	//         portNumber = self._get_vnc_fields(displayNumber)
	//         localPortNumber = find_free_port()

	//         vnc_tunnel = self._ssh.make_tunnel(localPortNumber, X_FWD_ADDR, portNumber)
	//         vnc_tunnel.start()
	//         vncviewer(otp, localPortNumber)
	//         self.refresh()

	return otp, displayNumber, nil
}

func (c *Connection) startVncSession(sizeX string, sizeY string) (string, string, error) {
	command := fmt.Sprintf("vncserver -geometry %sx%s -otp -novncauth -nohttpd", sizeX, sizeY)
	output, err := c.client.CleanExec(command)
	if err != nil {
		return "", "", nil
	}

	var otp, display string
	log.Println("Parsing vncserver stderr output.")

	for _, line := range strings.Split(output, "\n") {
		tline := strings.TrimSpace(line)
		//     # TurboVNC 1.1
		if strings.HasPrefix(tline, "New \x27X\x27 desktop is") {
			display = strings.Split(tline, " ")[4]
		}
		//     # TurboVNC 1.2
		if strings.HasPrefix(tline, "Desktop \x27TurboVNC:") {
			display = strings.Split(tline, " ")[2]
		}

		if strings.HasPrefix(tline, "Full control one-time password:") {
			otp = strings.Split(tline, " ")[4]
		}
	}
	displayNumber := strings.Split(display, ":")[1]

	return otp, displayNumber, nil
}

func (c *Connection) KillVncSession(display string) error {
	text, err := c.client.CleanExec(fmt.Sprintf("vncserver -kill %s", display))
	if err != nil {
		log.Println("Error killing vncserver; response: ", text)
	}
	return err
}

func parseSessions(username, text string) []Session {
	sessions := make([]Session, 0, 10)
	for _, line := range strings.Split(strings.TrimSuffix(text, "\n"), "\n") {
		session, err := parseVncLine(line)
		if err != nil {
			continue
		}
		if session.Username == username {
			sessions = append(sessions, session)
		}
	}
	return sessions
}

type Session struct {
	Username      string
	DisplayNumber string
	Geometry      string
	DateCreated   string
	HostPort      uint
}

func (s *Session) Label() string {
	return fmt.Sprintf(
		"Session %s - %s - %s",
		s.DisplayNumber,
		s.Geometry,
		s.DateCreated,
	)
}

func (s Session) AddListener(listener binding.DataListener) {
}

func (s Session) RemoveListener(listener binding.DataListener) {
}

func parseVncLine(line string) (Session, error) {
	var session Session
	fields := strings.Fields(line)
	if len(fields) < 15 {
		return session, errors.New("Parse error")
	}
	username := fields[15][1 : len(fields[15])-1]
	port, err := strconv.Atoi(get_field(fields, "-rfbport"))
	if err != nil {
		return Session{}, err
	}
	session = Session{
		Username:      username,
		DisplayNumber: fields[11],
		Geometry:      get_field(fields, "-geometry"),
		DateCreated:   fields[8],
		HostPort:      uint(port),
	}
	return session, nil
}

func get_field(fields []string, fieldname string) string {
	for ii, field := range fields {
		if field == fieldname {
			return fields[ii+1]
		}
	}
	return ""
}

type Tunnel struct {
	local_port   uint
	client       Client
	endpoint     string
	listener     net.Listener
	closeChannel chan interface{}
}

func (t Tunnel) LocalPort() uint {
	return t.local_port
}

func (t Tunnel) Close() {
	t.closeChannel <- 1
}

func acceptConnection(listener net.Listener, c chan net.Conn) {
	conn, err := listener.Accept()
	if err != nil {
		return
	}
	c <- conn
}

func (capcon *Connection) NewTunnel(local_p uint, remote_h string, remote_p uint) (*Tunnel, error) {
	endpoint := fmt.Sprintf("localhost:%d", local_p)
	listener, err := net.Listen("tcp", endpoint)
	if err != nil {
		return nil, err
	}
	remote_endpoint := fmt.Sprintf("%s:%d", remote_h, remote_p)
	tunnel := &Tunnel{
		local_p,
		capcon.client,
		remote_endpoint,
		listener,
		make(chan interface{}),
	}
	go tunnel.Start()
	return tunnel, nil
}

func (t Tunnel) Start() {
	defer t.listener.Close()
	for {
		c := make(chan net.Conn)
		go acceptConnection(t.listener, c)
		log.Println("listening for new connections...")
		select {

		case <-t.closeChannel:
			return

		case localConn := <-c:
			log.Println("accepted connection")
			go func() {
				remoteConn, err := t.client.Dial("tcp", t.endpoint)
				if err != nil {
					log.Printf("server dial error: %s", err)
					return
				}
				copyConn := func(writer, reader net.Conn) {
					_, err := io.Copy(writer, reader)
					if err != nil {
						log.Printf(
							"Error: could not copy %s -> %s because: %s",
							reader,
							writer,
							err,
						)
					}
				}
				go copyConn(localConn, remoteConn)
				go copyConn(remoteConn, localConn)
			}()

		}
	}

}
