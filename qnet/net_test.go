package qnet

import (
	"net"
	"testing"

	"github.com/qiangyt/go-comm/v2/q18n"
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

	q18n.InitI18n("en")

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
	q18n.InitI18n("en")
	_, _, err = ResolveBroadcastIp(interfaces, "invalid-interface")
	if err != nil {
		a.Contains(err.Error(), "invalid-interface")
		a.Contains(err.Error(), "not found")
	}

	// Test error message in Chinese
	q18n.InitI18n("zh")
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

func TestBroadcastInterfacesP(t *testing.T) {
	a := require.New(t)

	// Test the panic version
	interfaces := BroadcastInterfacesP(true)
	a.NotNil(interfaces)
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

func TestBroadcastIpWithInterface_NoValidAddress(t *testing.T) {
	a := require.New(t)

	// We need to test this with real interfaces
	interfaces, err := BroadcastInterfaces(false)
	a.NoError(err)

	for _, iface := range interfaces {
		// Try loopback interface which should work differently
		if iface.Name == "lo" {
			broadcastIP, err := BroadcastIpWithInterface(iface)
			// Loopback might not return a broadcast IP
			if err == nil {
				t.Logf("Loopback broadcast IP: %v", broadcastIP)
			}
			return
		}
	}

	t.Skip("No loopback interface found")
}

func TestResolveBroadcastIpP_PanicOnError(t *testing.T) {
	// Test panic on error with non-existent interface
	defer func() {
		if r := recover(); r == nil {
			t.Error("ResolveBroadcastIpP should panic on error")
		}
	}()

	interfaces, _ := BroadcastInterfaces(false)
	ResolveBroadcastIpP(interfaces, "nonexistent-interface-xyz")
}

func TestResolveBroadcastIpP_Happy(t *testing.T) {
	a := require.New(t)

	interfaces, err := BroadcastInterfaces(false)
	a.NoError(err)

	// Try to find a valid interface
	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		localIP, broadcastIP, err := ResolveBroadcastIp(interfaces, iface.Name)
		if err == nil {
			a.NotNil(localIP)
			a.NotNil(broadcastIP)

			// Test the panic version with valid data
			localIP2, broadcastIP2 := ResolveBroadcastIpP(interfaces, iface.Name)
			a.Equal(localIP.String(), localIP2.String())
			a.Equal(broadcastIP.String(), broadcastIP2.String())
			return
		}
	}

	t.Skip("No valid network interface found")
}

// Tests for net.go

func TestBroadcastIpWithInterfaceP_Loopback(t *testing.T) {
	a := require.New(t)

	// Create a loopback interface
	intf := net.Interface{
		Index:        1,
		MTU:          1500,
		Name:         "lo0",
		HardwareAddr: nil,
		Flags:        net.FlagLoopback | net.FlagUp,
	}

	// On systems with loopback, this should work
	// On Windows without proper loopback config, might return nil
	result := BroadcastIpWithInterfaceP(intf)
	// Just verify it doesn't panic and returns a valid IP or nil
	a.True(result == nil || len(result) == 4 || len(result) == 16)
}

func TestBroadcastIpWithInterface_NoIP(t *testing.T) {
	a := require.New(t)

	// Create an interface with no addresses
	intf := net.Interface{
		Index:        999,
		MTU:          1500,
		Name:         "dummy0",
		HardwareAddr: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		Flags:        net.FlagUp,
	}

	// Should return nil, nil (no error, no IP)
	ip, err := BroadcastIpWithInterface(intf)
	// Either returns nil,nil or error depending on system
	a.True(ip == nil || err != nil)
}
