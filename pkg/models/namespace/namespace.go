package namespace

import (
	"github.com/automated-pen-testing/api/pkg/models/project"
	"github.com/automated-pen-testing/api/pkg/models/user"

	"gorm.io/gorm"
)

type (
	// Namespace manage projects admin can create namespaces
	Namespace struct {
		gorm.Model
		Name     string
		Users    []*user.User       `gorm:"-"`
		Projects []*project.Project `gorm:"foreignKey:namespace_id"`
	}
)
