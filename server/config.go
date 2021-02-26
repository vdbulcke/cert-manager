package server

// Config config for server
type Config struct {
	ListeningAddress string `yaml:"address" validate:"required"`
	ListeningPort    string `yaml:"port" validate:"required"`
}
