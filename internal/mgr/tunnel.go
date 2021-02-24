package mgr

import (
	"context"
	"errors"

	"sshtunnel/internal/config"
	"sshtunnel/pkg/tunnel"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

var (
	// ErrTunnelNotExist is returned when can not find
	// the given tunnel in system
	ErrTunnelNotExist = errors.New("tunnel not exist")
)

// TunnelManager is responsible for managing tunnels
// This implementation only takes accountability for
// starting, terminating tunnels based on the provided configurations.
// There is no operations such as auto-starting failed tunnels, retries, backoff ...
type TunnelManager struct {
	config  *config.Config
	tunnels map[string]*tunnel.Tunnel
}

// New initializes a new TunnelManager
func New(config *config.Config) *TunnelManager {
	return &TunnelManager{
		config:  config,
		tunnels: make(map[string]*tunnel.Tunnel),
	}
}

// Start starts all tunnels defined in config
// If any tunnels fail to start it will continue starting other tunnels
// and logs all failed tunnels with error level
func (m *TunnelManager) Start() error {
	for _, t := range m.config.Tunnels {
		name, err := m.startTunnel(&t)
		if err != nil {
			log.Errorf("Failed to start tunnel: %s, error: %s", name, err)
		}
	}
	return nil
}

// Shutdown terminates all managed tunnels
// If any tunnels failed to shutdown we simply log it out
func (m *TunnelManager) Shutdown() error {
	for name := range m.tunnels {
		err := m.shutdownTunnel(name)
		if err != nil {
			log.Errorf("Failed to shutdown tunnel: %s, error: %s", name, err)
		}
	}
	return nil
}

// startTunnel starts a tunnel with the given tunnel configuration
func (m *TunnelManager) startTunnel(t *config.Tunnel) (string, error) {
	log.Infof("Starting tunnel: %s", t.Name)
	authMethod := ssh.Password(t.Auth.Password)
	tun := tunnel.New(t.Name, authMethod, t.Auth.Username, t.SSHAddr, t.RemoteAddr, t.LocalAddr)
	m.tunnels[t.Name] = tun
	go tun.Start()
	return t.Name, nil
}

// shutdownTunnel terminates a tunnel with the given name
func (m *TunnelManager) shutdownTunnel(name string) error {
	log.Infof("Shutting down tunnel: %s", name)
	tun, ok := m.tunnels[name]
	if !ok {
		return ErrTunnelNotExist
	}
	if tun == nil {
		return nil
	}
	err := tun.Shutdown(context.TODO())
	return err
}
