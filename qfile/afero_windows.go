//go:build windows
// +build windows

package qfile

func DefaultEtcHosts() (string, error) {
	return `C:\Windows\System32\Drivers\etc\hosts`, nil
}
