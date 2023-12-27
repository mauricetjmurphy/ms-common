package http

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultIdleConnTimeout       = 55 * time.Second
	defaultDialTimeout           = 300 * time.Second
	defaultDialKeepAlive         = 300 * time.Second
	defaultExpectContinueTimeout = 30 * time.Second
	defaultTLSHandshakeTimeout   = 5 * time.Second
)

type ClientConfig struct {
	// TransportConf provides the configurations for the HTTP default transport.
	TransportConf TransportConf
	// CheckRedirect specifies the policy for handling redirects.
	CheckRedirect func(req *http.Request, via []*http.Request) error
	// Optional TLS config for making secure requests
	TLSClientConfig *tls.Config
}

type TransportConf struct {
	MaxIdleConnsPerHost   int
	IdleConnTimeout       time.Duration
	TLSHandshakeTimeout   time.Duration
	ExpectContinueTimeout time.Duration
	DialContext           func(ctx context.Context, network, addr string) (net.Conn, error)
	Proxy                 func(*http.Request) (*url.URL, error)
	TLSNextProto          map[string]func(authority string, c *tls.Conn) http.RoundTripper
}

// NewDefaultClient builds default http.Client
func NewDefaultClient() *http.Client {
	return New(&ClientConfig{})
}

// New builds a new http.Client on given ClientConf.
// Usage:
//
//	httpClient := New(&http.ClientConf{})
func New(conf *ClientConfig) *http.Client {
	applyDefaults(conf)
	transport := newHTTPTransport(conf)
	return &http.Client{
		Transport:     transport,
		CheckRedirect: conf.CheckRedirect,
	}
}

func newHTTPTransport(conf *ClientConfig) *http.Transport {
	transport := conf.TransportConf
	return &http.Transport{
		DialContext:           transport.DialContext,
		ExpectContinueTimeout: transport.ExpectContinueTimeout,
		IdleConnTimeout:       transport.IdleConnTimeout,
		Proxy:                 transport.Proxy,
		MaxIdleConnsPerHost:   transport.MaxIdleConnsPerHost,
		TLSClientConfig:       conf.TLSClientConfig,
		TLSHandshakeTimeout:   transport.TLSHandshakeTimeout,
		TLSNextProto:          transport.TLSNextProto,
	}
}

func applyDefaults(conf *ClientConfig) {
	transport := &conf.TransportConf
	transport.ExpectContinueTimeout = getOrDefault(transport.ExpectContinueTimeout, defaultExpectContinueTimeout)
	transport.IdleConnTimeout = getOrDefault(transport.IdleConnTimeout, defaultIdleConnTimeout)
	transport.TLSHandshakeTimeout = getOrDefault(transport.TLSHandshakeTimeout, defaultTLSHandshakeTimeout)
	if transport.Proxy == nil {
		transport.Proxy = http.ProxyFromEnvironment
	}
	if conf.TransportConf.DialContext == nil {
		conf.TransportConf.DialContext = (&net.Dialer{
			Timeout:   defaultDialTimeout,
			KeepAlive: defaultDialKeepAlive,
		}).DialContext
	}
}

func getOrDefault(val, defaultVal time.Duration) time.Duration {
	switch {
	case val < 0:
		return time.Duration(0)
	case val == 0:
		return defaultVal
	default:
		return val
	}
}
