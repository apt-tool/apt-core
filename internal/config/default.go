package config

import (
	"github.com/apt-tool/apt-core/internal/config/core"
	"github.com/apt-tool/apt-core/internal/config/ftp"
	"github.com/apt-tool/apt-core/internal/config/migration"
	"github.com/apt-tool/apt-core/internal/storage/sql"
)

func Default() Config {
	return Config{
		Core: core.Config{
			Preemptive: false,
			Port:       8080,
			Enable:     false,
		},
		MySQL: sql.Config{
			Host:     "127.0.0.1",
			Port:     3306,
			User:     "root",
			Pass:     "",
			Database: "automated-pen-testing",
			Migrate:  false,
		},
		Migrate: migration.Config{
			Enable: false,
		},
		FTP: ftp.Config{
			Host:   "",
			Secret: "",
			Access: "",
		},
	}
}
