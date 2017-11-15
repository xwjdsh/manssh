package manssh

// HostConfig struct include alias, connect string and other config
type HostConfig struct {
	// Aliases may be multi, eg "a1 a2"
	Aliases string
	// Connect string format is user@host:port
	Connect string
	// Config is other configs
	Config map[string]string
}
