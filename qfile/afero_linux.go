//go:build linux
// +build linux

package qfile

func DefaultEtcHosts() (string, error) {
	return "/etc/hosts", nil
}
