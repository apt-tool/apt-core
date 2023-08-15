package project

import (
	"github.com/apt-tool/apt-core/pkg/models/document"

	"gorm.io/gorm"
)

type (
	Project struct {
		gorm.Model
		Name        string
		Description string
		Host        string
		Creator     string
		HTTPSecure  bool
		Port        int
		NamespaceID uint
		Endpoints   []*EndpointSet       `gorm:"foreignKey:project_id"`
		Labels      []*LabelSet          `gorm:"foreignKey:project_id"`
		Params      []*ParamSet          `gorm:"foreignKey:project_id"`
		Documents   []*document.Document `gorm:"foreignKey:project_id"`
	}

	EndpointSet struct {
		gorm.Model
		ProjectID uint
		Endpoint  string
	}

	ParamSet struct {
		gorm.Model
		ProjectID uint
		Key       string
		Value     string
	}

	LabelSet struct {
		gorm.Model
		ProjectID uint
		Key       string
		Value     string
	}
)
