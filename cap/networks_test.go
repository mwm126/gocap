package cap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetExternalIp(t *testing.T) {
	ip := GetExternalIp()

	assert.NotEqual(t, ip.String(), "127.0.0.1")
}
