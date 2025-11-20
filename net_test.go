package comm

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResolveBroadcastIp_Loopback(t *testing.T) {
	a := require.New(t)

	interfaces, err := BroadcastInterfaces(false)
	a.NoError(err)

	// Test with loopback interface (might not exist or not support broadcast)
	localIP, broadcastIP, err := ResolveBroadcastIp(interfaces, "lo")

	// On most systems, loopback interface doesn't support broadcast
	// So this might return an error
	if err != nil {
		a.Error(err)
	} else {
		a.NotNil(localIP)
		a.NotNil(broadcastIP)
	}
}

func TestResolveBroadcastIp_InvalidInterface(t *testing.T) {
	a := require.New(t)

	InitI18n("en")

	interfaces, err := BroadcastInterfaces(false)
	a.NoError(err)

	// Test with non-existent interface
	_, _, err = ResolveBroadcastIp(interfaces, "nonexistent-interface-xyz")

	a.Error(err)
	a.Contains(err.Error(), "nonexistent-interface-xyz")
}

func TestResolveBroadcastIp_I18n(t *testing.T) {
	a := require.New(t)

	interfaces, err := BroadcastInterfaces(false)
	a.NoError(err)

	// Test error message in English
	InitI18n("en")
	_, _, err = ResolveBroadcastIp(interfaces, "invalid-interface")
	if err != nil {
		a.Contains(err.Error(), "invalid-interface")
		a.Contains(err.Error(), "not found")
	}

	// Test error message in Chinese
	InitI18n("zh")
	_, _, err = ResolveBroadcastIp(interfaces, "invalid-interface")
	if err != nil {
		a.Contains(err.Error(), "invalid-interface")
		a.Contains(err.Error(), "接口")
	}
}

func TestBroadcastInterfaces(t *testing.T) {
	a := require.New(t)

	// Get all broadcast-capable interfaces
	interfaces, err := BroadcastInterfaces(false)
	a.NoError(err)
	a.NotNil(interfaces)

	// There should be at least some interfaces on the system
	t.Logf("Found %d broadcast-capable interfaces", len(interfaces))
}

func TestBroadcastIpWithInterface(t *testing.T) {
	a := require.New(t)

	// Get all interfaces
	allInterfaces, err := net.Interfaces()
	a.NoError(err)

	// Try to find at least one interface that supports broadcast
	foundValidInterface := false

	for _, iface := range allInterfaces {
		// Skip down interfaces
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		// Skip loopback
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// Check if it supports broadcast
		if iface.Flags&net.FlagBroadcast == 0 {
			continue
		}

		// Try to get broadcast IP for this interface
		broadcastIP, err := BroadcastIpWithInterface(iface)

		if err == nil {
			foundValidInterface = true
			a.NotNil(broadcastIP)
			a.Equal(byte(255), broadcastIP[len(broadcastIP)-1])
			t.Logf("Found valid broadcast interface: %s with IP: %v", iface.Name, broadcastIP)
			break
		}
	}

	// This test might not find a valid interface on all systems
	// So we don't fail if no valid interface is found
	if !foundValidInterface {
		t.Log("No valid network interface found for testing")
	}
}

func TestHostname(t *testing.T) {
	a := require.New(t)

	// Test getting hostname
	hostname, err := Hostname()
	a.NoError(err)
	a.NotEmpty(hostname)

	t.Logf("Hostname: %s", hostname)
}

func TestHostnameP(t *testing.T) {
	a := require.New(t)

	// Test getting hostname with panic version
	hostname := HostnameP()
	a.NotEmpty(hostname)
}

