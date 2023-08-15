package models

import (
	"github.com/apt-tool/apt-core/pkg/models/document"
	"github.com/apt-tool/apt-core/pkg/models/namespace"
	"github.com/apt-tool/apt-core/pkg/models/project"
	"github.com/apt-tool/apt-core/pkg/models/user"
	"github.com/apt-tool/apt-core/pkg/models/user_namespace"

	"gorm.io/gorm"
)

// Interface manages the models interfaces
type Interface struct {
	Documents     document.Interface
	Namespaces    namespace.Interface
	UserNamespace user_namespace.Interface
	Projects      project.Interface
	Users         user.Interface
}

func New(db *gorm.DB) *Interface {
	return &Interface{
		Documents:     document.New(db),
		Namespaces:    namespace.New(db),
		UserNamespace: user_namespace.New(db),
		Projects:      project.New(db),
		Users:         user.New(db),
	}
}
