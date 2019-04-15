package utils

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestIsPrivateIP(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.True(IsPrivateIP(net.ParseIP("127.0.0.1")))
	assert.False(IsPrivateIP(net.ParseIP("8.8.8.8")))
}

func TestIsPrivateHost(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.True(IsPrivateHost("localhost"))
	assert.True(IsPrivateHost("127.0.0.1"))
	assert.True(IsPrivateHost("192.168.2.1"))
	assert.False(IsPrivateHost("google.com"))
	assert.False(IsPrivateHost("trap.jp"))
	assert.False(IsPrivateHost("8.8.8.8"))
}
