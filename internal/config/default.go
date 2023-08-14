package config

import (
	"github.com/automated-pen-testing/api/internal/config/core"
	"github.com/automated-pen-testing/api/internal/config/ftp"
	"github.com/automated-pen-testing/api/internal/config/http"
	"github.com/automated-pen-testing/api/internal/config/migration"
	"github.com/automated-pen-testing/api/internal/storage/redis"
	"github.com/automated-pen-testing/api/internal/storage/sql"
	"github.com/automated-pen-testing/api/internal/utils/jwt"
)

func Default() Config {
	return Config{
		Core: core.Config{
			Preemptive: false,
			Port:       8080,
			Enable:     false,
		},
		HTTP: http.Config{
			Port: 8080,
			Core: "",
		},
		JWT: jwt.Config{
			PrivateKey: "private",
			ExpireTime: 60,
		},
		Redis: redis.Config{
			Host: "localhost:6379",
			Pass: "",
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
