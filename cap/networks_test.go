package cap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCapConnection(t *testing.T) {
	networks := getNetworks()

	assert.Equal(t, Network{}, networks["external"])
}
