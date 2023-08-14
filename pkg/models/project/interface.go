package project

import (
	"fmt"

	"gorm.io/gorm"
)

// Interface manages the project db methods
type Interface interface {
	Create(project *Project) error
	Delete(projectID uint) error
	GetByID(projectID uint) (*Project, error)
}

func New(db *gorm.DB) Interface {
	return &core{
		db: db,
	}
}

type core struct {
	db *gorm.DB
}

func (c core) Create(project *Project) error {
	return c.db.Create(project).Error
}

func (c core) Delete(projectID uint) error {
	return c.db.Delete(&Project{}, "id = ?", projectID).Error
}

func (c core) GetByID(projectID uint) (*Project, error) {
	project := new(Project)

	query := c.db.
		Preload("Documents").
		Preload("Labels").
		Preload("Endpoints").
		Preload("Params").
		First(&project, "id = ?", projectID)
	if err := query.Error; err != nil {
		return nil, fmt.Errorf("[db.Project.Get] failed to get record error=%w", err)
	}

	if project.ID != projectID {
		return nil, ErrProjectNotFound
	}

	return project, nil
}
