package cmd

import (
	"fmt"
	"log"

	"github.com/ptaas-tool/base-api/internal/config/migration"
	"github.com/ptaas-tool/base-api/internal/utils/crypto"
	"github.com/ptaas-tool/base-api/pkg/enum"
	"github.com/ptaas-tool/base-api/pkg/models/document"
	"github.com/ptaas-tool/base-api/pkg/models/namespace"
	"github.com/ptaas-tool/base-api/pkg/models/project"
	"github.com/ptaas-tool/base-api/pkg/models/user"
	"github.com/ptaas-tool/base-api/pkg/models/user_namespace"

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
