package cap

import (
	"bytes"
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

type sshClient struct {
	client *ssh.Client
}

func (sc sshClient) Close() {
	sc.client.Close()
}

func (sc sshClient) CleanExec(cmd string) (string, error) {
	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := sc.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	b, err := session.CombinedOutput(cmd)
	return string(b), err
}

func NewSshClient(server net.IP, user, pass string) (Client, error) {

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
		return nil, err
	}
	var sc sshClient
	sc.client = client

	return &sc, nil
}

func (sc *sshClient) OpenSSHTunnel(
	user, pass string,
	local_port int,
	remote_addr string,
	remote_port int,
) sshtunnel.SSHTunnel {
	//     Open forward to the login SSH daemon
	log.Println("Opening Tunnel")

	// ssh_port := check_free_port(SSH_LOCAL_PORT)

	tunnel := sshtunnel.NewSSHTunnel(
		sc.client,
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

func (client *sshClient) CheckPasswordExpired(
	pass string,
	request_password func(Client),
	ch chan string,
) error {
	if client.isPasswordExpired() {
		log.Println("Password expired.")
		request_password(client)
		new_password := <-ch
		log.Println("Got new password.")
		err := client.changePassword(pass, new_password)
		defer client.Close()
		if err != nil {
			log.Println("Unable to change password")
			return err
		}
		log.Println("Password changed.")
		return errors.New("Could not connect; password was expired")
	}
	return nil
}

func (client *sshClient) isPasswordExpired() bool {
	out, err := client.CleanExec("echo")
	if err != nil {
		log.Println("errTxt=", err)
	}
	log.Println("outTxt=", out)
	return strings.Contains(strings.ToLower(out), "expired")
}

func (client *sshClient) changePassword(
	old_pw string,
	newPasswd string,
) error {
	log.Println("Opening shell to existing connection")
	session, err := client.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		log.Println("Could not open xterm (pty)")
		return err
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		log.Println("Problem opening ssh stdin")
		return err
	}
	var buf bytes.Buffer
	session.Stdout = &buf
	session.Stderr = &buf

	err = session.Shell()
	if err != nil {
		log.Println("Problem opening ssh shell")
		return err
	}

	cmds := [][]string{
		{"Current Password", old_pw},
		{"New password", newPasswd},
		{"Retype new password", newPasswd},
	}

	for _, cmd := range cmds {
		expected_prompt := cmd[0]
		reply := cmd[1]

		for !strings.Contains(buf.String(), expected_prompt) {
			time.Sleep(10 * time.Second)
			log.Println("not found in string: ", buf.String()[len(buf.String())-20:])
			log.Println("Waiting for response: ", expected_prompt)
		}
		log.Println(">>>>>>>>>>>>>>>>>>>>>:  ", expected_prompt)
		_, err := fmt.Fprintf(stdin, "%s\n", reply)
		if err != nil {
			log.Println("Problem running command: ", err)
		}
		log.Println("<<<<<<<<<<<<<<<<<<<<<:  ", reply)
	}

	for !strings.Contains(buf.String(), "updated") {
		time.Sleep(1 * time.Second)
		log.Println("Expected: updated not found in string: ", buf.String())
	}
	return nil
}

func (client *sshClient) Dial(protocol, endpoint string) (net.Conn, error) {
	return client.client.Dial(protocol, endpoint)
}
