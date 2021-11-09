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

	"aeolustec.com/capclient/cap/sshtunnel"
	"golang.org/x/crypto/ssh"
)

type CapConnectionManager struct {
	knocker          Knocker
	connection       *CapConnection
	password_expired bool
}

func NewCapConnectionManager(knocker Knocker) *CapConnectionManager {
	return &CapConnectionManager{knocker, nil, false}
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
	pw_expired_cb func(PasswordChecker),
	ch chan string) error {
	log.Println("Opening SSH Connection...")
	err := cm.knocker.Knock(user, ext_addr, server)
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

func cleanExec(client *ssh.Client, cmd string) (string, error) {
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
	loginName, err := cleanExec(client, "hostname")
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
	connInfo := ConnectionInfo{
		user,
		pass,
		uid,
		loginName,
		loginAddr,
		webLocalPort,
		sshLocalPort.Local.Port,
	}
	var fwds map[string]sshtunnel.SSHTunnel
	return &CapConnection{client, connInfo, fwds}, nil
}

func getLoginIP(client *ssh.Client, loginName string) (string, error) {
	command := ("ping -c 1 " +
		loginName + "| grep PING|awk \x27{print $3}\x27" + "| sed \x22s/(//\x22|sed \x22s/)//\x22")
	log.Println("Command for getting address of login server: ", command)
	addr, err := cleanExec(client, command)
	if err != nil {
		return addr, err
	}
	return strings.TrimSpace(addr), nil
}

func getUID(client *ssh.Client) (string, error) {
	//  id|sed "s/uid=//"|sed "s/(/ /"|awk '{print $1}'
	command := "id|sed \x22s/uid=//\x22|sed \x22s/(/ /\x22" + "|awk \x27{print $1}\x27"
	log.Println("Command to get UID: ", command)
	uid, err := cleanExec(client, command)
	if err != nil {
		return uid, err
	}
	return strings.TrimSpace(uid), nil
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
