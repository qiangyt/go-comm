//go:build darwin
// +build darwin

package qio

func DefaultEtcHosts() (string, error) {
	return "/private/etc/hosts", nil
}
