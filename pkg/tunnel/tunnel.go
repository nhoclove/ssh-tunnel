package tunnel

import (
	"context"
	"io"
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// Tunnel represents a tunnel to forward traffic from local server
// to a remote server through SSH server
type Tunnel struct {
	name       string
	sshAddr    string
	remoteAddr string
	localAddr  string
	username   string
	password   string

	listener   *net.Listener
	onShutdown []func()
}

// New initializes a new tunnel
func New(name, sshAddr, remoteAddr, localAddr, username, password string) *Tunnel {
	return &Tunnel{
		name:       name,
		sshAddr:    sshAddr,
		remoteAddr: remoteAddr,
		localAddr:  localAddr,
		username:   username,
		password:   password,
	}
}

func initSSHConfig(user, pass string) (*ssh.ClientConfig, error) {
	config := ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return &config, nil
}

// Start starts a local server then forwards all connections to remote server
// through SSH tunnel server.
func (t *Tunnel) Start() error {
	local, err := net.Listen("tcp", t.localAddr)
	if err != nil {
		log.Fatalf("Failed to start local server: %v", err)
	}
	t.listener = &local
	for {
		conn, err := local.Accept()
		if err != nil {
			log.Errorf("Failed to accept a new connection: %v", err)
			return err
		}

		go t.forward(conn)
	}
}

// Shutdown terminates all servers
func (t *Tunnel) Shutdown(ctx context.Context) error {
	// Close listener
	var err error
	if lerr := (*t.listener).Close(); lerr != nil {
		log.Error("Failed to close listener: %v", lerr)
		err = lerr
	}
	return err
}

func (t *Tunnel) forward(localConn net.Conn) {
	log.Infof("New connection from: %s", localConn.LocalAddr())
	cfg, err := initSSHConfig(t.username, t.password)
	if err != nil {
		log.Fatalf("Failed to init SSH config: %v", err)
	}

	// Establish connection to SSH server
	sshConn, err := ssh.Dial("tcp", t.sshAddr, cfg)
	if err != nil {
		log.Fatalf("Failed to establish connection to SSH server: %v", err)
	}
	// Establish connection to remote server
	remoteConn, err := sshConn.Dial("tcp", t.remoteAddr)
	if err != nil {
		log.Fatalf("Failed to establish connection to Remote server: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go copy(localConn, remoteConn, &wg)
	wg.Add(1)
	go copy(remoteConn, localConn, &wg)
	wg.Wait()

	log.Infof("Close connection: %s", localConn.LocalAddr().String())
	sshConn.Close()
}

func copy(writer, reader net.Conn, wg *sync.WaitGroup) {
	defer func() {
		writer.Close()
		reader.Close()
		wg.Done()
	}()

	if _, err := io.Copy(writer, reader); err != nil {
		log.Errorf("Failed to copy io from %s to %s. Error: %s", reader.LocalAddr().String(), writer.LocalAddr().String(), err)
	}
}
