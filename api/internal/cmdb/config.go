package cmdb

// Config ...
type Config struct {
	Addr     string
	LogLevel string
	DSN      string
}

// NewConfig ...
func NewConfig(cfg Config) *Config {
	return &cfg
}
