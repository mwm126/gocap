package main

import (
	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

type SpyKnocker struct {
	username string
	password string
	network  string
	knocked  bool
}

func (sk *SpyKnocker) Knock(username, password, network string) {
	sk.knocked = true
	sk.username = username
	sk.password = password
	sk.network = network
}

func TestJouleLoginButton(t *testing.T) {
	spy := &SpyKnocker{}
	client := newClient(spy)
	test.Type(client.jouleTab.UsernameEntry, "the_user")
	test.Type(client.jouleTab.PasswordEntry, "the_pass")
	client.jouleTab.NetworkSelect.SetSelected("external")

	test.Tap(client.jouleTab.LoginBtn)

	assert.True(t, spy.knocked)
	assert.Equal(t, "the_user", spy.username)
	assert.Equal(t, "the_pass", spy.password)
	assert.Equal(t, "external", spy.network)
}

func TestWattLoginButton(t *testing.T) {
	spy := &SpyKnocker{}
	client := newClient(spy)
	test.Type(client.wattTab.UsernameEntry, "the_user")
	test.Type(client.wattTab.PasswordEntry, "the_pass")
	client.wattTab.NetworkSelect.SetSelected("external")

	test.Tap(client.wattTab.LoginBtn)

	assert.True(t, spy.knocked)
	assert.Equal(t, "the_user", spy.username)
	assert.Equal(t, "the_pass", spy.password)
	assert.Equal(t, "external", spy.network)
}
