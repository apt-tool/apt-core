package http

type Config struct {
	Port    int    `koanf:"port"`
	Core    string `koanf:"core"`
	DevMode bool   `koanf:"dev_mode"`
}
