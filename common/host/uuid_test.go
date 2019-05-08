package host

import (
	"testing"
)

func TestUUID(t *testing.T) {
	uuid := UUID()
	t.Log(uuid)
}

func TestIPAddress(t *testing.T) {
	ip := IPAddress()
	t.Log(ip)
}
