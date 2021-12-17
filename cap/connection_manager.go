package cap

import (
	"errors"
	"log"
	"net"
	"strings"

	"aeolustec.com/capclient/cap/sshtunnel"
)

type ConnectionManager struct {
	knocker          Knocker
	connection       *Connection
	password_expired bool
}

func NewCapConnectionManager(knocker Knocker) *ConnectionManager {
	return &ConnectionManager{knocker, nil, false}
}

func (t *ConnectionManager) Knocker() *Knocker {
	return &t.knocker
}

func (t *ConnectionManager) GetConnection() *Connection {
	return t.connection
}

func (c *ConnectionManager) GetPasswordExpired() bool {
	return c.password_expired
}

func (c *ConnectionManager) SetPasswordExpired() {
	c.password_expired = true
}

func (t *ConnectionManager) Close() {
	if t.connection == nil {
		log.Println("Not connected; Cannot close connection")
		return
	}
	t.connection.close()
	t.connection = nil
}

func (cm *ConnectionManager) Connect(
	user, pass string,
	ext_addr,
	server net.IP,
	port uint,
	pw_expired_cb func(Client),
	ch chan string) error {

	log.Println("Opening SSH Connection...")
	err := cm.knocker.Knock(user, ext_addr, server, port)
	if err != nil {
		return err
	}

	log.Println("Going to SSHClient.connect() to ", server, " with ", user)

	client, err := NewSshClient(server, user, pass)
	if err != nil {
		log.Println("Could not connect to", server, err)
		return err
	}

	// password_checker := PasswordChecker{client, pass}
	log.Println("Checking for expired password...")
	if client.IsPasswordExpired() {
		log.Println("Password expired.")
		pw_expired_cb(client)
		new_password := <-ch
		log.Println("Got new password.")
		err := client.ChangePassword(pass, new_password)
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

const webLocalPort = 10080

func NewCapConnection(client Client, user, pass string) (*Connection, error) {
	log.Println("Getting connection info...")
	loginName, err := client.CleanExec("hostname")
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

	sshLocalPort := client.OpenSSHTunnel(user, pass, SSH_LOCAL_PORT, SSH_FWD_ADDR, SSH_FWD_PORT)

	log.Println("Connected.")
	conn := Connection{
		client,
		make(map[string]sshtunnel.SSHTunnel),
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

func getLoginIP(client Client, loginName string) (string, error) {
	command := ("ping -c 1 " +
		loginName + "| grep PING|awk \x27{print $3}\x27" + "| sed \x22s/(//\x22|sed \x22s/)//\x22")
	log.Println("Command for getting address of login server: ", command)
	addr, err := client.CleanExec(command)
	if err != nil {
		return addr, err
	}
	return strings.TrimSpace(addr), nil
}

func getUID(client Client) (string, error) {
	//  id|sed "s/uid=//"|sed "s/(/ /"|awk '{print $1}'
	command := "id|sed \x22s/uid=//\x22|sed \x22s/(/ /\x22" + "|awk \x27{print $1}\x27"
	log.Println("Command to get UID: ", command)
	uid, err := client.CleanExec(command)
	if err != nil {
		return uid, err
	}
	return strings.TrimSpace(uid), nil
}

const SSH_LOCAL_PORT = 10022
const SSH_FWD_ADDR = "localhost"
const SSH_FWD_PORT = 22
