package sftp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Client presents functionalities to propagate via SFTP
//
//go:generate mockery --output sftpmocks --outpkg sftpmocks --name Client
type Client interface {
	// Put writes/uploads the given binary data to the target remote path.
	// It opens the named file with specified flags os.O_RDWR|os.O_CREATE|os.O_TRUNC
	// It returns the number of bytes written and an error, if any.
	//
	// Note that some SFTP servers (e.g. AWS Transfer) do not support opening files read/write at the same time.
	// For those services you will need to use specified flags os.O_WRONLY|os.O_CREATE|os.O_TRUNC
	Put(data []byte, path string, flag int) (int64, error)

	// Get downloads/opens named file with specified flags https://pkg.go.dev/os@go1.18#pkg-constants
	// It returns the bytes and an error, if any.
	Get(path string, flag int) ([]byte, error)

	// Close closes the SFTP session.
	Close() error
}

// clientImpl implements the Client interface.
type clientImpl struct {
	*sftp.Client
}

func (c *clientImpl) Put(data []byte, path string, flag int) (n int64, err error) {
	w, err := c.Client.OpenFile(path, flag)
	if err != nil {
		return 0, err
	}
	defer func(w *sftp.File) {
		err = w.Close()
	}(w)
	return io.Copy(w, bytes.NewReader(data))
}

func (c *clientImpl) Get(path string, flag int) (data []byte, err error) {
	r, err := c.Client.OpenFile(path, flag)
	if err != nil {
		return nil, err
	}
	defer func(r *sftp.File) {
		err = r.Close()
	}(r)
	return ioutil.ReadAll(r)
}

func (c *clientImpl) Close() error {
	return c.Client.Close()
}

// New creates the SFTP client instance wrapper sftp.Client
func New(client *sftp.Client) Client {
	return &clientImpl{client}
}

// Connect creates the SFTP client.
// It will be dialed connection the target SFTP on given configuration.
// It returns the Client and an error, if any.
func Connect(opts ...Opts) (Client, error) {
	cfg := &configs{
		port: DefaultPort,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	clientConfig, err := createClientConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "sftp: unable to load client config")
	}

	addr := cfg.GetAddr()
	conn, err := ssh.Dial(DialTCP, addr, clientConfig)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("sftp: unable to connect to %v", addr))
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}

	return &clientImpl{client}, nil
}

func createClientConfig(cfg *configs) (*ssh.ClientConfig, error) {
	var auths []ssh.AuthMethod
	if len(cfg.password) > 0 {
		auths = append(auths, ssh.Password(cfg.password))
	}
	if len(cfg.privateKey) > 0 {
		signer, err := parsePrivateKey(cfg)
		if err != nil {
			return nil, err
		}
		auths = append(auths, ssh.PublicKeys(signer))
	}

	var clientConfig = &ssh.ClientConfig{
		User: cfg.username,
		Auth: auths,
	}

	if cfg.hostKeyCallback != nil {
		clientConfig.HostKeyCallback = cfg.hostKeyCallback
	}

	if cfg.timeout != nil {
		clientConfig.Timeout = *cfg.timeout
	}

	return clientConfig, nil
}

func parsePrivateKey(cfg *configs) (ssh.Signer, error) {
	if len(cfg.keyPhrase) > 0 {
		singer, err := ssh.ParsePrivateKeyWithPassphrase([]byte(cfg.privateKey), []byte(cfg.keyPhrase))
		if err != nil {
			return nil, errors.Wrap(err, "sftp: unable to sign on private key with passphrase")
		}
		return singer, nil
	}
	signer, err := ssh.ParsePrivateKey([]byte(cfg.privateKey))
	if err != nil {
		return nil, errors.Wrap(err, "sftp: unable to sign on private key")
	}
	return signer, nil
}
