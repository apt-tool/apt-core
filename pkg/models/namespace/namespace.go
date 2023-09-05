package namespace

import (
	"github.com/apt-tool/apt-core/pkg/models/project"
	"github.com/apt-tool/apt-core/pkg/models/user"

	"gorm.io/gorm"
)

type (
	// Namespace manage projects admin can create namespaces
	Namespace struct {
		gorm.Model
		Name      string
		CreatedBy string
		Users     []*user.User       `gorm:"-"`
		Projects  []*project.Project `gorm:"foreignKey:namespace_id"`
	}
)
