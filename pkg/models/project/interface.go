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
	GetAll() ([]*Project, error)
}

func New(db *gorm.DB) Interface {
	return &core{
		db: db,
	}
}

type core struct {
	db *gorm.DB
}

// Create a new project
func (c core) Create(project *Project) error {
	return c.db.Create(project).Error
}

// Delete a existed project
func (c core) Delete(projectID uint) error {
	return c.db.Delete(&Project{}, "id = ?", projectID).Error
}

// GetByID returns a single project with all its dependencies
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

// GetAll projects without their dependencies
func (c core) GetAll() ([]*Project, error) {
	list := make([]*Project, 0)

	if err := c.db.Find(&list).Error; err != nil {
		return nil, fmt.Errorf("[db.Project.GetAll] failed to get records error=%w", err)
	}

	return list, nil
}
