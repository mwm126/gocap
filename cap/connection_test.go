// +build integration

package cap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCapConnection(t *testing.T) {
	log.Fatal("alkqdvkzjdkfj3")
	// username := "mmeredith"
	// password := "xUZv!jA&TgHTkw#!3$bUVcDXxW3sY"
	username := "mark"
	password := "markmark"
	server := "localhost:22"

	conn := newCapConnection(username, password, server)
	defer conn.close()

	guis := conn.listGUIs()

	assert.Equal(t, "/etc/passwd\n", guis)
}
