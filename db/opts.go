package db

// Option is a database configuration option.
type Option func(*Config)

// Host sets the remote hostname for the database to connect to.
func Host(host string) Option {
	return func(c *Config) {
		c.Host = host
	}
}

// Port sets the remote port for the database to connect to.
func Port(port int) Option {
	return func(c *Config) {
		c.Port = port
	}
}

// User sets the username for the database to use on connect.
func User(username string) Option {
	return func(c *Config) {
		c.User = username
	}
}

// Password sets the password to set on connect.
func Password(password string) Option {
	return func(c *Config) {
		c.Password = password
	}
}

// Name sets the name of the database to operate on.
func Name(databaseName string) Option {
	return func(c *Config) {
		c.Name = databaseName
	}
}

// AWSRegion sets the name of the AWS region to operate on.
func AWSRegion(region string) Option {
	return func(c *Config) {
		c.AWSRegion = region
	}
}

// AWSSecretID sets the name of the AWS secret ID to operate on.
func AWSSecretID(secretID string) Option {
	return func(c *Config) {
		c.AWSSecretID = secretID
	}
}

// MigrationSourceURL sets the name of the location path to store migration scripts to operate on.
func MigrationSourceURL(sourceFileURL string) Option {
	return func(c *Config) {
		c.MigrationSourceURL = sourceFileURL
	}
}
