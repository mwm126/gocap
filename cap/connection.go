package cap

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"aeolustec.com/capclient/client/sshtunnel"
	"golang.org/x/crypto/ssh"
)

type ConnectionManager interface {
	Connect(
		user, pass string,
		ext_addr,
		server net.IP,
		port uint,
		pw_expired_cb func(PasswordChecker),
		ch chan string) error
	Close()
	GetConnection() Connection
	GetPasswordExpired() bool
	SetPasswordExpired()
}

type CapConnectionManager struct {
	knocker          Knocker
	connection       *CapConnection
	password_expired bool
}

func NewCapConnectionManager(knocker Knocker) *CapConnectionManager {
	return &CapConnectionManager{knocker, nil, false}
}

func (t *CapConnectionManager) GetConnection() Connection {
	return t.connection
}

func (c *CapConnectionManager) GetPasswordExpired() bool {
	return c.password_expired
}

func (c *CapConnectionManager) SetPasswordExpired() {
	c.password_expired = true
}

func (t *CapConnectionManager) Close() {
	if t.connection == nil {
		log.Println("Not connected; Cannot close connection")
		return
	}
	t.connection.close()
	t.connection = nil
}

type Connection interface {
	FindSessions() ([]Session, error)
	CreateVncSession(xres string, yres string) (string, string, error)
	GetUsername() string
	UpdateForwards(fwds []string)
}

// A CapConnection represents a successful SSH connection after the port knock
type CapConnection struct {
	client       *ssh.Client
	forwards     map[string]sshtunnel.SSHTunnel
	username     string
	password     string
	uid          string
	loginName    string
	loginAddr    string
	webLocalPort int
	sshLocalPort int
}

func (c *CapConnection) GetUsername() string {
	return c.username
}

func (conn *CapConnection) UpdateForwards(fwds []string) {
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

func (conn *CapConnection) forward(fwd string) sshtunnel.SSHTunnel {

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
	return openSSHTunnel(conn.client, conn.username,
		conn.password, local_p, remote_h, remote_p)
}

func (conn *CapConnection) close() {
	for _, fwd := range conn.forwards {
		fwd.Close()
	}
	conn.client.Close()
}

func (cm *CapConnectionManager) Connect(
	user, pass string,
	ext_addr,
	server net.IP,
	port uint,
	pw_expired_cb func(PasswordChecker),
	ch chan string) error {

	log.Println("Opening SSH Connection...")
	err := cm.knocker.Knock(user, ext_addr, server, port)
	if err != nil {
		return err
	}

	log.Println("Going to SSHClient.connect() to ", server, " with ", user)

	// var hostKey ssh.PublicKey
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		// HostKeyCallback: ssh.FixedHostKey(hostKey),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	host := fmt.Sprintf("%s:%s", server, "22")
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		log.Println("Could not connect to", host, err)
		return err
	}

	password_checker := PasswordChecker{client, pass}
	log.Println("Checking for expired password...")
	if password_checker.is_pw_expired() {
		log.Println("Password expired.")
		pw_expired_cb(password_checker)
		new_password := <-ch
		log.Println("Got new password.")
		err := password_checker.change_password(client, pass, new_password)
		defer client.Close()
		if err != nil {
			log.Println("Unable to change password")
			return err
		}
		log.Println("Password changed.")
		return errors.New("Could not connect; password was expired")
	}

	cm.connection, err = NewCapConnection(client, user, pass)
	if err != nil {
		defer client.Close()
	}
	return err
}

func CleanExec(client *ssh.Client, cmd string) (string, error) {
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	b, err := session.CombinedOutput(cmd)
	return string(b), err
}

const webLocalPort = 10080

func NewCapConnection(client *ssh.Client, user, pass string) (*CapConnection, error) {
	log.Println("Getting connection info...")
	loginName, err := CleanExec(client, "hostname")
	if err != nil {
		log.Println("Failed hostname")
		return nil, err
	}
	loginName = strings.TrimSpace(loginName)
	log.Println("Got login hostname:", loginName)

	loginAddr, err := getLoginIP(client, loginName)
	if err != nil {
		log.Println("Failed to lookup login IP")
		return nil, err
	}
	log.Println("Got login server IP:", loginAddr)

	uid, err := getUID(client)
	if err != nil {
		log.Println("Failed to lookup UID")
		return nil, err
	}
	log.Println("Got UID:", uid)

	sshLocalPort := openSSHTunnel(client, user, pass, SSH_LOCAL_PORT, SSH_FWD_ADDR, SSH_FWD_PORT)

	log.Println("Connected.")
	conn := CapConnection{
		client,
		make(map[string]sshtunnel.SSHTunnel, 0),
		user,
		pass,
		uid,
		loginName,
		loginAddr,
		webLocalPort,
		sshLocalPort.Local.Port,
	}
	return &conn, nil
}

func getLoginIP(client *ssh.Client, loginName string) (string, error) {
	command := ("ping -c 1 " +
		loginName + "| grep PING|awk \x27{print $3}\x27" + "| sed \x22s/(//\x22|sed \x22s/)//\x22")
	log.Println("Command for getting address of login server: ", command)
	addr, err := CleanExec(client, command)
	if err != nil {
		return addr, err
	}
	return strings.TrimSpace(addr), nil
}

func getUID(client *ssh.Client) (string, error) {
	//  id|sed "s/uid=//"|sed "s/(/ /"|awk '{print $1}'
	command := "id|sed \x22s/uid=//\x22|sed \x22s/(/ /\x22" + "|awk \x27{print $1}\x27"
	log.Println("Command to get UID: ", command)
	uid, err := CleanExec(client, command)
	if err != nil {
		return uid, err
	}
	return strings.TrimSpace(uid), nil
}

const SSH_LOCAL_PORT = 10022
const SSH_FWD_ADDR = "localhost"
const SSH_FWD_PORT = 22

const VNC_LOCAL_PORT = 10055

func openSSHTunnel(
	client *ssh.Client,
	user, pass string,
	local_port int,
	remote_addr string,
	remote_port int,
) sshtunnel.SSHTunnel {
	//     Open forward to the login SSH daemon
	log.Println("Opening Tunnel")

	// ssh_port := check_free_port(SSH_LOCAL_PORT)

	tunnel := sshtunnel.NewSSHTunnel(
		client,
		// User and host of tunnel server, it will default to port 22
		// if not specified.
		fmt.Sprintf("%s@%s", user, "localhost"),

		// Pick ONE of the following authentication methods:
		// sshtunnel.PrivateKeyFile("path/to/private/key.pem"), // 1. private key
		ssh.Password(pass), // 2. password
		// sshtunnel.SSHAgent(),                                // 3. ssh-agent

		// The destination host and port of the actual server.
		fmt.Sprintf("%s:%d", remote_addr, remote_port),

		// The local port you want to bind the remote port to.
		// Specifying "0" will lead to a random port.
		strconv.Itoa(local_port),
	)

	tunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	go func() {
		err := tunnel.Start()
		if err != nil {
			log.Println("Could not create tunnel: ", err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
	log.Println("tunnel is ", tunnel)
	return *tunnel
}

// func readSessionManagerSecret(client *ssh.Client) string {
// 	command := "chmod 0400 ~/.sessionManager;cat ~/.sessionManager"
// 	out, err := CleanExec(client, command)

// 	if err == nil {
// 		log.Println("Reading .sessionManager secret")
// 		return out
// 	}

// 	log.Println("Error:  ", err)
// 	log.Println("Creating new .sessionManager secret")
// 	command = `dd if=/dev/urandom bs=1 count=1024|sha256sum|
//                awk \x27{print $1}\x27> ~/.sessionManager;
//                cat ~/.sessionManager;chmod 0400 ~/.sessionManager`
// 	out, err = CleanExec(client, command)
// 	return out
// }

func (c *CapConnection) FindSessions() ([]Session, error) {
	sessions := make([]Session, 0, 10)
	text, err := CleanExec(c.client, "ps auxnww|grep Xvnc|grep -v grep")
	if err != nil {
		return sessions, errors.New("Unable to find sessions")
	}
	return parseSessions(c.GetUsername(), text), nil
}

const MAX_GUI_COUNT = 4

func (c *CapConnection) CreateVncSession(xres string, yres string) (string, string, error) {
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

func (c *CapConnection) startVncSession(sizeX string, sizeY string) (string, string) {
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
