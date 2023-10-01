package config

import (
	"github.com/ptaas-tool/base-api/internal/config/core"
	"github.com/ptaas-tool/base-api/internal/config/ftp"
	"github.com/ptaas-tool/base-api/internal/config/migration"
	"github.com/ptaas-tool/base-api/internal/storage/sql"
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
