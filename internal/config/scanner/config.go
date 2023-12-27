package scanner

type Config struct {
	Command  string   `koanf:"command"`
	Enable   bool     `koanf:"enable"`
	Flags    []string `koanf:"flags"`
	Defaults []string `koanf:"defaults"`
}
