//go:build darwin
// +build darwin

package qfile

func DefaultEtcHosts() (string, error) {
	return "/private/etc/hosts", nil
}
