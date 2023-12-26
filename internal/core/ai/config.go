package ai

type Config struct {
	Method string `koanf:"method"`
	Factor int    `koanf:"factor"`
	Limit  int    `koanf:"limit"`
	Enable bool   `koanf:"enable"`
}
