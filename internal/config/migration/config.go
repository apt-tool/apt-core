package migration

type Config struct {
	Root   string `koanf:"root"`
	Pass   string `koanf:"pass"`
	Enable bool   `koanf:"enable"`
}
