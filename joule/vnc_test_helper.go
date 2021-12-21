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

func NewFakeVncConnection(t *testing.T) *cap.Connection {
	conn, err := cap.NewCapConnection(&FakeVncClient{}, "the_user", "pass")
	if err != nil {
		t.Error(err)
	}
	return conn
}

func NewFakeVncClient(server net.IP, user, pass string) (cap.Client, error) {
	client := FakeVncClient{}
	return &client, nil
}

type FakeVncClient struct{}

func (fsc FakeVncClient) CleanExec(command string) (string, error) {
	replies := map[string]string{
		"hostname": "the_hostname",
		`ping -c 1 the_hostname| grep PING|awk '{print $3}'| sed "s/(//"|sed "s/)//"`: "1.2.3.4",
		`id|sed "s/uid=//"|sed "s/(/ /"|awk '{print $1}'`:                             "the_uid",
		"ps auxnww|grep Xvnc|grep -v grep":                                            `8227 27248  0.0  0.0  92620 70572 ?        S    Aug03   0:234 /nfs/apps/TurboVNC/2.0.2/bin/Xvnc :123 -desktop TurboVNC: login03:5 (the_user) -auth /nfs/home/3/mmeredith/.Xauthority -dontdisconnect -geometry 3840x2160 -depth 24 -rfbwait 120000 -otpauth -pamauth -rfbport 5905 -fp catalogue:/etc/X11/fontpath.d -deferupdate 1`,
	}
	reply, exists := replies[command]
	if !exists {
		log.Println(command, " : ", reply)
		return "", errors.New("Unexpected command")
	}
	return replies[command], nil
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
