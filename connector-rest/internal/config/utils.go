package config

func AssembleAddress(config Config) string {
	return config.HTTPServer.ListeningAddress + config.HTTPServer.Port
}
