package mgr

import (
	"errors"
	"io/ioutil"
	"sshtunnel/internal/config"
	"strings"

	"golang.org/x/crypto/ssh"
)

func (m *TunnelManager) getAuthMethods(auth config.Auth) ([]ssh.AuthMethod, error) {
	if strings.TrimSpace(auth.User) == "" {
		return nil, errors.New("empty auth user")
	}
	methods := make([]ssh.AuthMethod, 0, 2)
	keyPath := strings.TrimSpace(auth.KeyPath)
	if keyPath != "" {
		buff, err := ioutil.ReadFile(keyPath)
		if err != nil {
			return nil, err
		}
		key, err := ssh.ParsePrivateKey(buff)
		if err != nil {
			return nil, err
		}
		methods = append(methods, ssh.PublicKeys(key))
	}
	password := strings.TrimSpace(auth.Password)
	if password != "" {
		methods = append(methods, ssh.Password(password))
	}
	if len(methods) < 1 {
		return nil, errors.New("no auth method provided")
	}
	return methods, nil
}
