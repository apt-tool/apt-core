package core

type Config struct {
	Timeout    int    `koanf:"timeout"`
	Port       int    `koanf:"port"`
	Workers    int    `koanf:"workers"`
	Enable     bool   `koanf:"enable"`
	Preemptive bool   `koanf:"preemptive"`
	Secret     string `koanf:"secret"`
}
