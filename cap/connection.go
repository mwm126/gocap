package cap

import (
	"fmt"
	"github.com/elliotchance/sshtunnel"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"time"
)

// A CapConnection represents a successful SSH connection after the port knock
type CapConnection struct {
	client         *ssh.Client
	session        *ssh.Session
	connectionInfo ConnectionInfo
}

// ConnectionInfo stores info about the connection
type ConnectionInfo struct {
	username string
	password string
	// _login_info,
	uid          string
	loginName    string
	loginAddr    string
	webLocalPort int
	sshLocalPort int
	mgtPort      sshtunnel.SSHTunnel
}

func (conn *CapConnection) listGUIs() string {
	out, err := conn.session.CombinedOutput("ls /etc/passwd")
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func (conn *CapConnection) close() {
	conn.client.Close()
}

func newCapConnection(user, pass, server string, knckr Knocker) (*CapConnection, error) {
	log.Println("Sending CAP packet...")
	knckr.Knock(user, pass, server)
	log.Println("Opening SSH Connection...")

	//     self.ssh.set_missing_host_key_policy(paramiko.AutoAddPolicy())
	log.Println("Going to SSHClient.connect() to ", server, " with ", user)
	client, session, err := connectToHost(user, pass, "localhost:22")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	//     password_checker = PasswordChecker(self.ssh, self._login_info.passwd)
	log.Println("Checking for expired password...")
	//     if password_checker.is_pw_expired():
	//         self.status_signal.emit("Password expired.")
	//         return (ConnectionEvent.GET_NEW_PASSWORD.value, password_checker)

	//     return (self._conn_event.value, self.get_connection())
	conn := getConnection(client, session, user, pass)
	return &conn, nil
}

func connectToHost(user, pass, host string) (*ssh.Client, *ssh.Session, error) {

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		log.Fatal(err)
		return nil, nil, err
	}

	return client, session, nil
}

const webLocalPort = 10080

func getConnection(client *ssh.Client, session *ssh.Session, user, pass string) CapConnection {

	log.Println("Getting connection info...")
	loginName := getHostname()
	loginAddr := getLoginIP(loginName)
	uid := getUID()
	sshLocalPort := openSSHTunnel(user, pass)

	log.Println("Starting session manager...")
	session_mgt_port := openSessionManagementForward(user, pass)
	//     session_mgt_secret = self._readSessionManagerSecret()

	log.Println("Connected.")
	//     sess_man = SessionManager(self._config, uid, session_mgt_port, session_mgt_secret)
	connInfo := ConnectionInfo{
		user,
		pass,
		// _login_info,
		uid,
		loginName,
		loginAddr,
		webLocalPort,
		sshLocalPort.Local.Port,
		session_mgt_port,
	}
	//     return HpcConnection(self.ssh, conn_info, sess_man)
	conn := CapConnection{client, session, connInfo}
	return conn
}

func getHostname() string {
	log.Println("Getting hostname")
	return "hostname"
}

func getLoginIP(loginName string) string {
	//     Get the login IP Address
	log.Println("Getting login IP")
	//     command = (
	//         f"ping -c 1 {loginName}"
	//         "| grep PING|awk \x27{print $3}\x27"
	//         "| sed \x22s/(//\x22|sed \x22s/)//\x22"
	//     )
	return "DO PING"
}

func getUID() string {
	//     Get uid for user
	//     # id|sed "s/uid=//"|sed "s/(/ /"|awk '{print $1}'
	log.Println("Getting login UID")
	// command := "id|sed \x22s/uid=//\x22|sed \x22s/(/ /\x22" + "|awk \x27{print $1}\x27"
	return "THE_UID"
}

const SSH_LOCAL_PORT = "10022"
const SSH_FWD_ADDR = "localhost"
const SSH_FWD_PORT = 22

func openSSHTunnel(user, pass string) sshtunnel.SSHTunnel {
	//     Open forward to the login SSH daemon
	log.Println("Opening Tunnel")

	// ssh_port := check_free_port(SSH_LOCAL_PORT)
	ssh_port := SSH_LOCAL_PORT

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

func openSessionManagementForward(user, pass string) sshtunnel.SSHTunnel {
	log.Println("Connecting to session manager")

	// ssh_port := check_free_port(SSH_LOCAL_PORT)
	ssh_port := SSH_LOCAL_PORT

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

func readSessionManagerSecret(conn CapConnection) []byte {
	command := "chmod 0400 ~/.sessionManager;cat ~/.sessionManager"
	out, err := conn.session.CombinedOutput(command)

	if err == nil {
		log.Println("Reading .sessionManager secret")
		return out
	}

	log.Println("Error:  ", err)
	log.Println("Creating new .sessionManager secret")
	command = `dd if=/dev/urandom bs=1 count=1024|sha256sum|
               awk \x27{print $1}\x27> ~/.sessionManager;
               cat ~/.sessionManager;chmod 0400 ~/.sessionManager`
	out, err = conn.session.CombinedOutput(command)
	return out
}
