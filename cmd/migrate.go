package cmd

import (
	"fmt"
	"log"

	"github.com/apt-tool/apt-core/internal/config/migration"
	"github.com/apt-tool/apt-core/internal/utils/crypto"
	"github.com/apt-tool/apt-core/pkg/enum"
	"github.com/apt-tool/apt-core/pkg/models/document"
	"github.com/apt-tool/apt-core/pkg/models/namespace"
	"github.com/apt-tool/apt-core/pkg/models/project"
	"github.com/apt-tool/apt-core/pkg/models/user"
	"github.com/apt-tool/apt-core/pkg/models/user_namespace"

	"gorm.io/gorm"
)

// Migrate is the command of migration
type Migrate struct {
	Cfg migration.Config
	Db  *gorm.DB
}

func (m Migrate) Do() {
	models := []interface{}{
		&document.Document{},
		&namespace.Namespace{},
		&user_namespace.UserNamespace{},
		&project.ParamSet{},
		&project.LabelSet{},
		&project.EndpointSet{},
		&project.Project{},
		&user.User{},
	}

	for _, item := range models {
		if err := m.Db.AutoMigrate(item); err != nil {
			log.Println(fmt.Errorf("failed to migrate model error=%w", err))
		}
	}

	if m.Cfg.Enable {
		tmp := &user.User{
			Username: m.Cfg.Root,
			Password: crypto.GetMD5Hash(m.Cfg.Pass),
			Role:     enum.RoleAdmin,
		}

		if err := m.Db.Create(tmp).Error; err != nil {
			log.Println(fmt.Errorf("failed to insert root user error=%w", err))
		}

		log.Println("root created!")
	}
}
