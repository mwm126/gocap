package cap

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"aeolustec.com/capclient/cap/sshtunnel"
	"golang.org/x/crypto/ssh"
)

type CapConnectionManager struct {
	knocker    Knocker
	connection *CapConnection
}

func NewCapConnectionManager(knocker Knocker) *CapConnectionManager {
	return &CapConnectionManager{knocker, nil}
}

func (t *CapConnectionManager) GetConnection() *CapConnection {
	return t.connection
}

func (t *CapConnectionManager) Close() {
	if t.connection == nil {
		log.Println("Not connected; Cannot close connection")
		return
	}
	t.connection.close()
	t.connection = nil
}

// A CapConnection represents a successful SSH connection after the port knock
type CapConnection struct {
	client         *ssh.Client
	connectionInfo ConnectionInfo
	forwards       map[string]sshtunnel.SSHTunnel
}

// ConnectionInfo stores info about the connection
type ConnectionInfo struct {
	username     string
	password     string
	uid          string
	loginName    string
	loginAddr    string
	webLocalPort int
	sshLocalPort int
	mgtPort      sshtunnel.SSHTunnel
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
	return openSSHTunnel(conn.client, conn.connectionInfo.username,
		conn.connectionInfo.password, local_p, remote_h, remote_p)
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
) error {
	log.Println("Opening SSH Connection...")
	err := cm.knocker.Knock(user, ext_addr, server)
	if err != nil {
		return err
	}

	//     self.ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
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
		log.Println("Could not connect to localhost:22, ", err)
		return err
	}

	//     password_checker = PasswordChecker(self.ssh, self._login_info.passwd)
	// log.Println("Checking for expired password...")
	//     if password_checker.is_pw_expired():
	//         self.status_signal.emit("Password expired.")
	//         return (ConnectionEvent.GET_NEW_PASSWORD.value, password_checker)

	//     return (self._conn_event.value, self.get_connection())
	err = cm.setupConnection(client, user, pass)
	if err != nil {
		client.Close()
	}
	return err
}

func cleanExec(client *ssh.Client, cmd string) (string, error) {
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	b, err := session.Output(cmd)
	return string(b), err
}

const webLocalPort = 10080

func (cm *CapConnectionManager) setupConnection(client *ssh.Client, user, pass string) error {
	log.Println("Getting connection info...")
	loginName, err := cleanExec(client, "hostname")
	if err != nil {
		log.Println("Failed hostname")
		return err
	}
	loginAddr, err := getLoginIP(client, loginName)
	if err != nil {
		log.Println("Failed to lookup login IP")
		return err
	}
	uid, err := getUID(client)
	if err != nil {
		log.Println("Failed to lookup UID")
		return err
	}
	sshLocalPort := openSSHTunnel(client, user, pass, SSH_LOCAL_PORT, SSH_FWD_ADDR, SSH_FWD_PORT)

	session_mgt_port := openSessionManagementForward(client, user, pass)
	// session_mgt_secret := readSessionManagerSecret(client)

	log.Println("Connected.")
	//     sess_man = SessionManager(self._config, uid, session_mgt_port, session_mgt_secret)
	connInfo := ConnectionInfo{
		user,
		pass,
		uid,
		loginName,
		loginAddr,
		webLocalPort,
		sshLocalPort.Local.Port,
		session_mgt_port,
	}
	var fwds map[string]sshtunnel.SSHTunnel
	cm.connection = &CapConnection{client, connInfo, fwds}
	return nil
}

func getLoginIP(client *ssh.Client, loginName string) (string, error) {
	//     Get the login IP Address
	log.Println("Getting login IP")
	command := ("ping -c 1 " +
		loginName + "| grep PING|awk \x27{print $3}\x27" + "| sed \x22s/(//\x22|sed \x22s/)//\x22")
	return cleanExec(client, command)
}

func getUID(client *ssh.Client) (string, error) {
	//     Get uid for user
	//     # id|sed "s/uid=//"|sed "s/(/ /"|awk '{print $1}'
	log.Println("Getting login UID")
	command := "id|sed \x22s/uid=//\x22|sed \x22s/(/ /\x22" + "|awk \x27{print $1}\x27"
	return cleanExec(client, command)
}

const SSH_LOCAL_PORT = 10022
const SSH_FWD_ADDR = "localhost"
const SSH_FWD_PORT = 22

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

func openSessionManagementForward(client *ssh.Client, user, pass string) sshtunnel.SSHTunnel {
	log.Println("Starting session manager...")

	// ssh_port := check_free_port(SSH_LOCAL_PORT)
	ssh_port := strconv.Itoa(SSH_LOCAL_PORT)

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
		fmt.Sprintf("%s:%d", SSH_FWD_ADDR, SSH_FWD_PORT),

		// The local port you want to bind the remote port to.
		// Specifying "0" will lead to a random port.
		ssh_port,
	)

	tunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
	go func() {
		err := tunnel.Start()
		if err != nil {
			log.Println("Could not create session management tunnel: ", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)
	log.Println("tunnel is ", tunnel)
	return *tunnel
}

// func readSessionManagerSecret(client *ssh.Client) string {
// 	command := "chmod 0400 ~/.sessionManager;cat ~/.sessionManager"
// 	out, err := cleanExec(client, command)

// 	if err == nil {
// 		log.Println("Reading .sessionManager secret")
// 		return out
// 	}

// 	log.Println("Error:  ", err)
// 	log.Println("Creating new .sessionManager secret")
// 	command = `dd if=/dev/urandom bs=1 count=1024|sha256sum|
//                awk \x27{print $1}\x27> ~/.sessionManager;
//                cat ~/.sessionManager;chmod 0400 ~/.sessionManager`
// 	out, err = cleanExec(client, command)
// 	return out
// }
