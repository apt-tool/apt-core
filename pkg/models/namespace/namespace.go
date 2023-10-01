package namespace

import (
	"github.com/ptaas-tool/base-api/pkg/models/project"
	"github.com/ptaas-tool/base-api/pkg/models/user"

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
