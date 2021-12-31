package cap

import (
	"log"
	"net"
	"strings"

	"aeolustec.com/capclient/cap/sshtunnel"
)

type ClientFactory func(server net.IP, user, pass, port string) (Client, error)

type ConnectionManager struct {
	clientFactory      ClientFactory
	knocker            *Knocker
	password_expired   bool
	NewPasswordChannel chan string
}

func NewCapConnectionManager(cf ClientFactory, knocker *Knocker) *ConnectionManager {
	return &ConnectionManager{cf, knocker, false, make(chan string)}
}

func (t *ConnectionManager) AddYubikeyCallback(cb func(bool)) {
	t.knocker.AddCallback(cb)
}

func (c *ConnectionManager) GetPasswordExpired() bool {
	return c.password_expired
}

func (cm *ConnectionManager) Connect(
	user, pass string,
	ext_addr,
	server net.IP,
	cap_port uint,
	ssh_port string,
	request_password func(Client),
	ch chan string) (*Connection, error) {

	log.Println("Opening SSH Connection...", cm.knocker.Yubikey)
	err := cm.knocker.Knock(user, ext_addr, server, cap_port)
	if err != nil {
		return nil, err
	}

	log.Println("Going to SSHClient.connect() to ", server, " with ", user)

	client, err := cm.clientFactory(server, user, pass, ssh_port)
	if err != nil {
		log.Println("Could not connect to", server, err)
		return nil, err
	}

	log.Println("Checking for expired password...")
	if err := client.CheckPasswordExpired(pass, request_password, ch); err != nil {
		cm.password_expired = true
		return nil, err
	}

	conn, err := NewCapConnection(client, user, pass)
	if err != nil {
		defer client.Close()
		return nil, err
	}

	return conn, nil
}

const webLocalPort = 10080

func NewCapConnection(client Client, user, pass string) (*Connection, error) {
	hostName, err := client.CleanExec("hostname")
	if err != nil {
		log.Println("Failed hostname")
		return nil, err
	}
	hostName = strings.TrimSpace(hostName)
	log.Println("Got login hostname:", hostName)

	loginAddr, err := getLoginIP(client, hostName)
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
	return &Connection{
		client,
		make(map[string]sshtunnel.SSHTunnel),
		user,
		pass,
		uid,
		hostName,
		loginAddr,
		webLocalPort,
		sshLocalPort.Local.Port,
	}, nil
}

func getLoginIP(client Client, hostName string) (string, error) {
	command := ("ping -c 1 " +
		hostName + "| grep PING|awk \x27{print $3}\x27" + "| sed \x22s/(//\x22|sed \x22s/)//\x22")
	addr, err := client.CleanExec(command)
	if err != nil {
		return addr, err
	}
	return strings.TrimSpace(addr), nil
}

func getUID(client Client) (string, error) {
	//  id|sed "s/uid=//"|sed "s/(/ /"|awk '{print $1}'
	command := "id|sed \x22s/uid=//\x22|sed \x22s/(/ /\x22" + "|awk \x27{print $1}\x27"
	uid, err := client.CleanExec(command)
	if err != nil {
		return uid, err
	}
	return strings.TrimSpace(uid), nil
}

const SSH_LOCAL_PORT = 10022
const SSH_FWD_ADDR = "localhost"
const SSH_FWD_PORT = 22
