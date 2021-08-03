package cap

import (
	"golang.org/x/crypto/ssh"
	"log"
)

// A CapConnection represents a successful SSH connection after the port knock
type CapConnection struct {
	client         *ssh.Client
	session        *ssh.Session
	connectionInfo ConnectionInfo
}

// ConnectionInfo stores info about the connection
type ConnectionInfo struct {
	// _login_info,
	uid          string
	loginName    string
	loginAddr    string
	webLocalPort int
	sshLocalPort string
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
		return nil, err
	}
	return nil, err

	//     password_checker = PasswordChecker(self.ssh, self._login_info.passwd)
	log.Println("Checking for expired password...")
	//     if password_checker.is_pw_expired():
	//         self.status_signal.emit("Password expired.")
	//         return (ConnectionEvent.GET_NEW_PASSWORD.value, password_checker)

	//     return (self._conn_event.value, self.get_connection())
	conn := getConnection(client, session)
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
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	return client, session, nil
}

const webLocalPort = 10080

func getConnection(client *ssh.Client, session *ssh.Session) CapConnection {

	log.Println("Getting connection info...")
	loginName := getHostname()
	loginAddr := getLoginIP(loginName)
	uid := getUID()
	sshLocalPort := openSSHTunnel()

	log.Println("Starting session manager...")
	//     session_mgt_port = self._openSessionManagementForward()
	//     session_mgt_secret = self._readSessionManagerSecret()

	log.Println("Connected.")
	//     sess_man = SessionManager(self._config, uid, session_mgt_port, session_mgt_secret)
	connInfo := ConnectionInfo{
		// _login_info,
		uid,
		loginName,
		loginAddr,
		webLocalPort,
		sshLocalPort,
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

func openSSHTunnel() string {
	//     Open forward to the login SSH daemon
	log.Println("Opening Tunnel")
	//     ssh_port = check_free_port(SSH_LOCAL_PORT)
	//     self._ssh_tunnel = Thread(
	//         target=tunnelTCP,
	//         args=(ssh_port, SSH_FWD_ADDR, SSH_FWD_PORT, self.ssh.get_transport()),
	//     )
	//     self._ssh_tunnel.daemon = True
	//     self._ssh_tunnel.start()
	return "ssh_port"
}

func openSessionManagementForward() int {
	//     Open forward to the session management server
	log.Println("Connecting to session manager")
	//     with closing(socket.socket(socket.AF_INET, socket.SOCK_STREAM)) as s:
	//         s.bind(("", 0))
	//         session_mgt_port = s.getsockname()[1]

	//     self._session_mgt_tunnel = Thread(
	//         target=tunnelTCP,
	//         args=(
	//             session_mgt_port,
	//             SESSION_MGT_FWD_ADDR,
	//             SESSION_MGT_FWD_PORT,
	//             self.ssh.get_transport(),
	//         ),
	//     )
	//     self._session_mgt_tunnel.daemon = True
	//     self._session_mgt_tunnel.start()
	//     return session_mgt_port
	return 1234
}

func readSessionManagerSecret() []byte {
	//     Read session manager secret
	//     command = "chmod 0400 ~/.sessionManager;cat ~/.sessionManager"
	//     _, stdout, stderr = self._cleanExec(command)

	//     errorOut = stderr.read()
	//     if errorOut == "":
	//         LOG.debug("Reading .sessionManager secret")
	//         return stdout.read()

	//     LOG.debug("Error:  %s", repr(errorOut))
	//     LOG.debug("Creating new .sessionManager secret")
	//     command = (
	//         "dd if=/dev/urandom bs=1 count=1024|sha256sum|"
	//         "awk \x27{print $1}\x27> ~/.sessionManager;"
	//         "cat ~/.sessionManager;chmod 0400 ~/.sessionManager"
	//     )
	//     _, stdout, _ = self._cleanExec(command)
	//     return stdout.read()
	return []byte{1, 2, 3, 4, 5}
}
