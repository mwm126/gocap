package main

import (
	"fmt"
	"github.com/elliotchance/sshtunnel"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"time"
)

func maingo() {
	// Setup the tunnel, but do not yet start it yet.
	tunnel := sshtunnel.NewSSHTunnel(
		// User and host of tunnel server, it will default to port 22
		// if not specified.
		"mark@localhost",

		// Pick ONE of the following authentication methods:
		// sshtunnel.PrivateKeyFile("path/to/private/key.pem"), // 1. private key
		ssh.Password("markmark"), // 2. password
		// sshtunnel.SSHAgent(),                                // 3. ssh-agent

		// The destination host and port of the actual server.
		"localhost:22",

		// The local port you want to bind the remote port to.
		// Specifying "0" will lead to a random port.
		"2222",
	)

	// You can provide a logger for debugging, or remove this line to
	// make it silent.
	tunnel.Log = log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)

	// Start the server in the background. You will need to wait a
	// small amount of time for it to bind to the localhost port
	// before you can start sending connections.
	go tunnel.Start()
	time.Sleep(100 * time.Millisecond * 1000)

	// NewSSHTunnel will bind to a random port so that you can have
	// multiple SSH tunnels available. The port is available through:
	//   tunnel.Local.Port

	// You can use any normal Go code to connect to the destination server
	// through localhost. You may need to use 127.0.0.1 for some libraries.
	//
	// Here is an example of connecting to a PostgreSQL server:
	// conn := fmt.Sprintf("host=127.0.0.1 port=%d username=foo", tunnel.Local.Port)
	// db, _ := sql.Open("postgres", conn)
	fmt.Println("ENDING.")

	// ...

}
