package sftp

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	DialTCP     = "tcp"
	DefaultPort = 22
)

type configs struct {
	host            string
	port            uint
	username        string
	password        string
	privateKey      string
	keyPhrase       string
	hostKeyCallback ssh.HostKeyCallback
	timeout         *time.Duration
}

type Opts func(*configs)

func (c *configs) GetAddr() string {
	return net.JoinHostPort(c.host, fmt.Sprint(c.port))
}

// Host sets SFTP hostname
func Host(hostname string) Opts {
	return func(c *configs) {
		c.host = hostname
	}
}

// Port sets SFTP port. Default 22
func Port(port uint) Opts {
	return func(c *configs) {
		c.port = port
	}
}

// Username sets SFTP username
func Username(username string) Opts {
	return func(c *configs) {
		c.username = username
	}
}

// Password sets SFTP secret password
func Password(password string) Opts {
	return func(c *configs) {
		c.password = password
	}
}

// PrivateKey sets SFTP private key
func PrivateKey(privateKey string) Opts {
	return func(c *configs) {
		c.privateKey = privateKey
	}
}

// KeyPhrase sets SFTP passphrase password
func KeyPhrase(keyPhrase string) Opts {
	return func(c *configs) {
		c.keyPhrase = keyPhrase
	}
}

// HostKeyCallback sets SFTP host key call function
func HostKeyCallback(hostKeyCallback ssh.HostKeyCallback) Opts {
	return func(c *configs) {
		c.hostKeyCallback = hostKeyCallback
	}
}

// Timeout is the maximum amount of time for the TCP connection to establish
func Timeout(timeout *time.Duration) Opts {
	return func(c *configs) {
		c.timeout = timeout
	}
}
