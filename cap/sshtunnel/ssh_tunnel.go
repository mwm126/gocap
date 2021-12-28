package sshtunnel

import (
	"io"
	"net"

	"golang.org/x/crypto/ssh"
)

type logger interface {
	Printf(string, ...interface{})
}

type SSHTunnel struct {
	Client   *ssh.Client
	Local    *Endpoint
	Server   *Endpoint
	Remote   *Endpoint
	Log      logger
	Conns    []net.Conn
	SvrConns []*net.Conn
	isOpen   bool
	close    chan interface{}
}

func (tunnel *SSHTunnel) logf(fmt string, args ...interface{}) {
	if tunnel.Log != nil {
		tunnel.Log.Printf(fmt, args...)
	}
}

func newConnectionWaiter(listener net.Listener, c chan net.Conn) {
	conn, err := listener.Accept()
	if err != nil {
		return
	}
	c <- conn
}

func (tunnel *SSHTunnel) Start() error {
	listener, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		return err
	}
	tunnel.isOpen = true
	tunnel.Local.Port = listener.Addr().(*net.TCPAddr).Port

	for {
		if !tunnel.isOpen {
			break
		}

		c := make(chan net.Conn)
		go newConnectionWaiter(listener, c)
		tunnel.logf("listening for new connections...")

		select {
		case <-tunnel.close:
			tunnel.logf("close signal received, closing...")
			tunnel.isOpen = false
		case conn := <-c:
			tunnel.Conns = append(tunnel.Conns, conn)
			tunnel.logf("accepted connection")
			go tunnel.forward(conn)
		}
	}
	total := len(tunnel.Conns)
	for i, conn := range tunnel.Conns {
		tunnel.logf("closing the netConn (%d of %d)", i+1, total)
		err := conn.Close()
		if err != nil {
			tunnel.logf(err.Error())
		}
	}
	total = len(tunnel.SvrConns)
	for i, conn := range tunnel.SvrConns {
		tunnel.logf("closing the serverConn (%d of %d)", i+1, total)
		err := (*conn).Close()
		if err != nil {
			tunnel.logf(err.Error())
		}
	}
	err = listener.Close()
	if err != nil {
		return err
	}
	tunnel.logf("tunnel closed")
	return nil
}

func (tunnel *SSHTunnel) forward(localConn net.Conn) {
	serverConn, err := tunnel.Client.Dial("tcp", tunnel.Server.String())
	if err != nil {
		tunnel.logf("server dial error: %s", err)
		return
	}
	tunnel.logf("connected to %s (1 of 2)\n", tunnel.Server.String())
	tunnel.SvrConns = append(tunnel.SvrConns, &serverConn)

	remoteConn, err := serverConn, nil
	if err != nil {
		tunnel.logf("remote dial error: %s", err)
		return
	}
	tunnel.Conns = append(tunnel.Conns, remoteConn)
	tunnel.logf("connected to %s (2 of 2)\n", tunnel.Remote.String())
	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			tunnel.logf("io.Copy error: %s", err)
		}
	}
	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

func (tunnel *SSHTunnel) Close() {
	tunnel.close <- struct{}{}
}

// NewSSHTunnel creates a new single-use tunnel. Supplying "0" for localport will use a random port.
func NewSSHTunnel(
	client *ssh.Client,
	tunnel string,
	auth ssh.AuthMethod,
	destination string,
	localport string,
) *SSHTunnel {

	localEndpoint := NewEndpoint("localhost:" + localport)

	server := NewEndpoint(tunnel)
	if server.Port == 0 {
		server.Port = 22
	}

	sshTunnel := &SSHTunnel{
		Local:  localEndpoint,
		Client: client,
		Server: server,
		Remote: NewEndpoint(destination),
		close:  make(chan interface{}),
	}

	return sshTunnel
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
