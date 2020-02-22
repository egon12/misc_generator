package map_graphql_to_proto

// Config variable to store global config
type Config struct {
	RemovePrefixes []string
	ProtoPackage   string
}

// Global config
var config Config

// SetConfig set global config
func SetConfig(c Config) {
	config = c
}

// GetConfig get global config
func GetConfig() Config {
	return config
}
