//go:build linux
// +build linux

package qio

func DefaultEtcHosts() (string, error) {
	return "/etc/hosts", nil
}
