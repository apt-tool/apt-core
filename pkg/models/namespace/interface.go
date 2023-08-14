package namespace

import (
	"fmt"

	"gorm.io/gorm"
)

// Interface manages the namespace db methods
type Interface interface {
	Create(namespace *Namespace) error
	Delete(namespaceID uint) error
	Get() ([]*Namespace, error)
	GetByID(namespaceID uint) (*Namespace, error)
	GetByIDs(namespaceIDs []uint) ([]*Namespace, error)
}

func New(db *gorm.DB) Interface {
	return &core{
		db: db,
	}
}

type core struct {
	db *gorm.DB
}

func (c core) Create(namespace *Namespace) error {
	return c.db.Create(namespace).Error
}

func (c core) Delete(namespaceID uint) error {
	return c.db.Delete(&Namespace{}, "id = ?", namespaceID).Error
}

func (c core) Get() ([]*Namespace, error) {
	list := make([]*Namespace, 0)

	if err := c.db.Find(&list).Error; err != nil {
		return nil, fmt.Errorf("[db.Namespace.Get] failed to get records error=%w", err)
	}

	return list, nil
}

func (c core) GetByID(namespaceID uint) (*Namespace, error) {
	namespace := new(Namespace)

	if err := c.db.Preload("Projects").Where("id = ?", namespaceID).First(&namespace).Error; err != nil {
		return nil, fmt.Errorf("[db.Namespace.GetByID] failed to get record error=%w", err)
	}

	if namespace.ID != namespaceID {
		return nil, ErrRecordNotFound
	}

	return namespace, nil
}

func (c core) GetByIDs(namespaceIDs []uint) ([]*Namespace, error) {
	list := make([]*Namespace, 0)

	if err := c.db.Where("id in ?", namespaceIDs).Find(&list).Error; err != nil {
		return nil, fmt.Errorf("[db.Namespace.Get] failed to get records error=%w", err)
	}

	return list, nil
}
