package infrastructure

// Listener configuration object
type Listener struct {
	Network string `toml:"network"`
	Address string `toml:"address"`
}

// Mongo configuration object
type Mongo struct {
	URL        string `toml:"url"`
	DB         string `toml:"db"`
	Collection string `toml:"connection"`
}

// ServerConfig config for MongoDB and network listener
type ServerConfig struct {
	Listener *Listener
	Mongo    *Mongo
}

// NewServerConfig creates empty server config
func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Listener: &Listener{},
		Mongo:    &Mongo{},
	}
}
