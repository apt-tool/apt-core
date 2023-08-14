package redis

type Config struct {
	Host string `koanf:"host"`
	Pass string `koanf:"pass"`
}
