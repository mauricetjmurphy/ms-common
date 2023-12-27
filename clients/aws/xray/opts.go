package xray

type config struct {
	enable     bool
	daemonAddr string
	version    string
}

type Option func(*config)

func WithEnable(enable bool) func(*config) {
	return func(c *config) {
		c.enable = enable
	}
}

func WithHost(daemonAddr string) func(*config) {
	return func(c *config) {
		c.daemonAddr = daemonAddr
	}
}

func WithVersion(version string) func(*config) {
	return func(c *config) {
		c.version = version
	}
}
