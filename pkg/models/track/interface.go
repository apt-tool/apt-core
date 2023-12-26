package track

import (
	"fmt"

	"gorm.io/gorm"
)

// Interface manages the documents methods
type Interface interface {
	Create(track *Track) error
	Get(id uint, projectID uint) ([]*Track, error)
}

func New(db *gorm.DB) Interface {
	return &core{
		db: db,
	}
}

type core struct {
	db *gorm.DB
}

// Create new track
func (c core) Create(track *Track) error {
	return c.db.Create(track).Error
}

// Get all tracks by project id
func (c core) Get(id uint, projectID uint) ([]*Track, error) {
	tracks := make([]*Track, 0)

	query := c.db.Model(&Track{}).Where("id > ?", id).Where("project_id = ?", projectID)
	if err := query.Find(&tracks).Error; err != nil {
		return nil, fmt.Errorf("[db.Document.GetByID] failed to get record error=%w", err)
	}

	return tracks, nil
}
