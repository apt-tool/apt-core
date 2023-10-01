package document

import (
	"fmt"

	"gorm.io/gorm"
)

// Interface manages the documents methods
type Interface interface {
	Create(document *Document) error
	Update(document *Document) error
	GetByID(id uint) (*Document, error)
}

func New(db *gorm.DB) Interface {
	return &core{
		db: db,
	}
}

type core struct {
	db *gorm.DB
}

func (c core) Create(document *Document) error {
	return c.db.Create(document).Error
}

func (c core) Update(document *Document) error {
	return c.db.Save(document).Error
}

func (c core) GetByID(id uint) (*Document, error) {
	document := new(Document)

	query := c.db.First(&document, "id = ?", id)
	if err := query.Error; err != nil {
		return nil, fmt.Errorf("[db.Document.GetByID] failed to get record error=%w", err)
	}

	if document.ID != id {
		return nil, ErrDocumentNotFound
	}

	return document, nil
}
