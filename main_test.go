package main

import (
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SpyKnocker struct {
	username string
	password string
	knocked  bool
}

func (sk *SpyKnocker) Knock(username, password string) {
	sk.knocked = true
}

func aTestLoginButton(t *testing.T) {
	spy := &SpyKnocker{}
	client := newClient(spy)
	test.Type(client.usernameEntry, "the_user")
	test.Type(client.passwordEntry, "the_pass")

	test.Tap(client.loginBtn)

	assert.True(t, spy.knocked)
	assert.Equal(t, "the_user", spy.username)
	assert.Equal(t, "the_pass", spy.password)
}
