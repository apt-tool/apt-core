package ftp

type Config struct {
	Host   string `koanf:"host"`
	Secret string `koanf:"secret"`
	Access string `koanf:"access"`
}
