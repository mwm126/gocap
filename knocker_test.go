package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPacket(t *testing.T) {
	pk := &PortKnocker{}
	pkt := pk.packet()

	assert.Equal(t, "TODO", pkt)
}
