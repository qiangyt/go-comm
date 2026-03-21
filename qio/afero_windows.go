//go:build windows
// +build windows

package qio

func DefaultEtcHosts() (string, error) {
	return `C:\Windows\System32\Drivers\etc\hosts`, nil
}
