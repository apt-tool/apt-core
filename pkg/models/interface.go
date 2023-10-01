package models

import (
	"github.com/ptaas-tool/base-api/pkg/models/document"
	"github.com/ptaas-tool/base-api/pkg/models/project"
	"github.com/ptaas-tool/base-api/pkg/models/user"

	"gorm.io/gorm"
)

// Interface manages the models interfaces
type Interface struct {
	Documents document.Interface
	Projects  project.Interface
	Users     user.Interface
}

func New(db *gorm.DB) *Interface {
	return &Interface{
		Documents: document.New(db),
		Projects:  project.New(db),
		Users:     user.New(db),
	}
}
