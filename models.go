package manssh

// HostConfig struct include alias, connect string and other config
type HostConfig struct {
	// Aliases may be multi, eg "a1 a2"
	Aliases string
	// Connect string format is user@host:port
	Connect string
	// Path found in which file
	Path string
	// Config is other configs
	Config map[string]string
}
