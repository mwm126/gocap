package joule

import (
	"errors"
	"io"
	"log"
	"net"
	"testing"

	"aeolustec.com/capclient/cap"
	"aeolustec.com/capclient/cap/sshtunnel"
	"golang.org/x/crypto/ssh"
)

func NewFakeVncConnection(t *testing.T, expectations map[string]string) *cap.Connection {
	conn, err := cap.NewCapConnection(&FakeVncClient{t, expectations}, "the_user", "pass")
	if err != nil || conn == nil {
		t.Fatal(err)
	}
	return conn
}

type FakeVncClient struct {
	t            *testing.T
	expectations map[string]string
}

func (fsc FakeVncClient) CleanExec(command string) (string, error) {
	reply, exists := fsc.expectations[command]
	if !exists {
		log.Println(command, " : ", reply)
		fsc.t.Fatal("Could not find response for command: ", command)
		return "", errors.New("Not reachable")
	}
	return fsc.expectations[command], nil
}

func (fsc FakeVncClient) Close() {
}

func (client FakeVncClient) CheckPasswordExpired(
	pass string,
	pw_expired_cb func(cap.Client),
	ch chan string,
) error {
	return nil
}

func (sc FakeVncClient) OpenSSHTunnel(
	user, pass string,
	local_port int,
	remote_addr string,
	remote_port int,
) sshtunnel.SSHTunnel {
	return *sshtunnel.NewSSHTunnel(
		nil,
		"testuser@localhost",
		ssh.Password(pass),
		"rem_addr:123",
		"123",
	)
}

func (sc FakeVncClient) Dial(protocol, endpoint string) (net.Conn, error) {
	return nil, nil
}

func (fsc *FakeVncClient) Start(command string) (io.ReadCloser, io.ReadCloser, error) {
	return nil, nil, nil
}

func (fsc *FakeVncClient) Wait() error {
	return nil
}

type StubYubikey struct{}

func (yk *StubYubikey) FindSerial() (int32, error) {
	return 0, nil
}

func (yk *StubYubikey) ChallengeResponse(chal [6]byte) ([16]byte, error) {
	return [16]byte{}, nil
}

func (yk *StubYubikey) ChallengeResponseHMAC(chal cap.SHADigest) ([20]byte, error) {
	return [20]byte{}, nil
}
