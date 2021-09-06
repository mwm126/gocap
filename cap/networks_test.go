package cap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetExternalIp(t *testing.T) {
	ip := GetExternalIp()

	assert.NotEqual(t, ip.String(), "127.0.0.1")
}
