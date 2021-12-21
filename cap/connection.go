package cap

import (
	"errors"
	"fmt"
	"log"
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
}

// A Connection represents a successful SSH connection after the port knock
type Connection struct {
	client       Client
	forwards     map[string]sshtunnel.SSHTunnel
	username     string
	password     string
	uid          string
	loginName    string
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
	c.client.Close()
}

// func readSessionManagerSecret(client Client) string {
//	command := "chmod 0400 ~/.sessionManager;cat ~/.sessionManager"
//	out, err := CleanExec(client, command)

//	if err == nil {
//		log.Println("Reading .sessionManager secret")
//		return out
//	}

//	log.Println("Error:  ", err)
//	log.Println("Creating new .sessionManager secret")
//	command = `dd if=/dev/urandom bs=1 count=1024|sha256sum|
//                awk \x27{print $1}\x27> ~/.sessionManager;
//                cat ~/.sessionManager;chmod 0400 ~/.sessionManager`
//	out, err = CleanExec(client, command)
//	return out
// }

func (c *Connection) FindSessions() ([]Session, error) {
	sessions := make([]Session, 0, 10)
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
		log.Println(tline)
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
	HostAddress   string
	HostPort      string
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
	fields := strings.Fields(line)
	username := fields[15][1 : len(fields[15])-1]
	session := Session{
		Username:      username,
		DisplayNumber: fields[11],
		Geometry:      get_field(fields, "-geometry"),
		DateCreated:   fields[8],
		HostAddress:   "localhost",
		HostPort:      get_field(fields, "-rfbport"),
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
