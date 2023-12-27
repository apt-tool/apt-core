package scanner

type Config struct {
	Command  string   `koanf:"command"`
	Enable   bool     `koanf:"enable"`
	Defaults []string `koanf:"defaults"`
}
