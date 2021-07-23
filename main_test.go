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

func (sk *SpyKnocker) Knock(uname string, pword string) {
	sk.knocked = true
	sk.username = uname
	sk.password = pword
}

func TestLoginButton(t *testing.T) {
	spy := &SpyKnocker{}
	client := newClient(spy)
	test.Type(client.username_entry, "the_user")
	test.Type(client.password_entry, "the_pass")

	test.Tap(client.login_btn)

	assert.True(t, spy.knocked)
	assert.Equal(t, "the_user", spy.username)
	assert.Equal(t, "the_pass", spy.password)
}
