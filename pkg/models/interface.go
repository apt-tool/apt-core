package models

import (
	"github.com/ptaas-tool/base-api/pkg/models/document"
	"github.com/ptaas-tool/base-api/pkg/models/namespace"
	"github.com/ptaas-tool/base-api/pkg/models/project"
	"github.com/ptaas-tool/base-api/pkg/models/user"
	"github.com/ptaas-tool/base-api/pkg/models/user_namespace"

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
