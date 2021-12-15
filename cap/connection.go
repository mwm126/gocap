package cap

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"aeolustec.com/capclient/cap/sshtunnel"
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

	otp, displayNumber := c.startVncSession(xres, yres)

	//         portNumber = self._get_vnc_fields(displayNumber)
	//         localPortNumber = find_free_port()

	//         vnc_tunnel = self._ssh.make_tunnel(localPortNumber, X_FWD_ADDR, portNumber)
	//         vnc_tunnel.start()
	//         vncviewer(otp, localPortNumber)
	//         self.refresh()

	return otp, displayNumber, nil
}

func (c *Connection) startVncSession(sizeX string, sizeY string) (string, string) {
	// command = f"vncserver -geometry {sizeX}x{sizeY} -otp -novncauth -nohttpd"
	// LOG.debug(command)

	// _stdin, _stdout, stderr = self._ssh.cleanExec(command)

	// display = ""
	otp := ""
	// LOG.debug("Parsing vncserver stderr output.")
	// for line in stderr:
	//     LOG.debug(line)
	//     # TurboVNC 1.1
	//     if line.strip().startswith("New \x27X\x27 desktop is"):
	//         display = str(line.split()[4].rstrip())
	//     # TurboVNC 1.2
	//     if line.strip().startswith("Desktop \x27TurboVNC:"):
	//         display = str(line.split()[2].rstrip())

	//     if line.strip().startswith("Full control one-time password:"):
	//         otp = str(line.split()[4].rstrip())
	displayNumber := ""
	// foundColon = False
	// for c in display:
	//     if foundColon:
	//         displayNumber += c
	//     if c == ":":
	//         foundColon = True

	return otp, displayNumber
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
