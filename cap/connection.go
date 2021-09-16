package cap

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/elliotchance/sshtunnel"
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

func (t *CapConnectionManager) CloseConnection() {
	if t.connection == nil {
		log.Println("Cannot close connection (connection missing)")
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
	return openSSHTunnel(conn.connectionInfo.username,
		conn.connectionInfo.password, local_p, remote_h, remote_p)
}

func (conn *CapConnection) close() {
	for _, fwd := range conn.forwards {
		fwd.Close()
	}
	conn.client.Close()
}

func (cm *CapConnectionManager) newCapConnection(
	user, pass string,
	server net.IP,
) (*CapConnection, error) {
	log.Println("Opening SSH Connection...")
	cm.knocker.Knock(user, pass, server)

	//     self.ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
	log.Println("Going to SSHClient.connect() to ", server, " with ", user)
	client, err := connectToHost(user, pass, "localhost:22")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	//     password_checker = PasswordChecker(self.ssh, self._login_info.passwd)
	log.Println("Checking for expired password...")
	//     if password_checker.is_pw_expired():
	//         self.status_signal.emit("Password expired.")
	//         return (ConnectionEvent.GET_NEW_PASSWORD.value, password_checker)

	//     return (self._conn_event.value, self.get_connection())
	cm.connection = getConnection(client, user, pass)
	return cm.connection, nil
}

func connectToHost(user, pass, host string) (*ssh.Client, error) {
	// var hostKey ssh.PublicKey
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		// HostKeyCallback: ssh.FixedHostKey(hostKey),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		log.Println("Failed to dial: ", err)
		return nil, err
	}
	// defer client.Close()
	return client, nil
}

func cleanExec(client *ssh.Client, cmd string) (string, error) {
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(cmd); err != nil {
		log.Fatal("Failed to run: " + err.Error())
		return "", err
	}
	return b.String(), nil
}

const webLocalPort = 10080

func getConnection(client *ssh.Client, user, pass string) *CapConnection {
	log.Println("Getting connection info...")
	loginName, err := cleanExec(client, "hostname")
	if err != nil {
		log.Fatal("Failed hostname")
	}
	loginAddr, err := getLoginIP(client, loginName)
	if err != nil {
		log.Fatal("Failed to lookup login IP")
	}
	uid, err := getUID(client)
	sshLocalPort := openSSHTunnel(user, pass, SSH_LOCAL_PORT, SSH_FWD_ADDR, SSH_FWD_PORT)

	session_mgt_port := openSessionManagementForward(user, pass)
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
	return &CapConnection{client, connInfo, fwds}
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
	user, pass string,
	local_port int,
	remote_addr string,
	remote_port int,
) sshtunnel.SSHTunnel {
	//     Open forward to the login SSH daemon
	log.Println("Opening Tunnel")

	// ssh_port := check_free_port(SSH_LOCAL_PORT)

	tunnel := sshtunnel.NewSSHTunnel(
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
	go tunnel.Start()
	time.Sleep(100 * time.Millisecond)
	log.Println("tunnel is ", tunnel)
	return *tunnel
}

func openSessionManagementForward(user, pass string) sshtunnel.SSHTunnel {
	log.Println("Starting session manager...")

	// ssh_port := check_free_port(SSH_LOCAL_PORT)
	ssh_port := strconv.Itoa(SSH_LOCAL_PORT)

	tunnel := sshtunnel.NewSSHTunnel(
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
	go tunnel.Start()
	time.Sleep(100 * time.Millisecond)
	log.Println("tunnel is ", tunnel)
	return *tunnel
}

func readSessionManagerSecret(client *ssh.Client) string {
	command := "chmod 0400 ~/.sessionManager;cat ~/.sessionManager"
	out, err := cleanExec(client, command)

	if err == nil {
		log.Println("Reading .sessionManager secret")
		return out
	}

	log.Println("Error:  ", err)
	log.Println("Creating new .sessionManager secret")
	command = `dd if=/dev/urandom bs=1 count=1024|sha256sum|
               awk \x27{print $1}\x27> ~/.sessionManager;
               cat ~/.sessionManager;chmod 0400 ~/.sessionManager`
	out, err = cleanExec(client, command)
	return out
}
