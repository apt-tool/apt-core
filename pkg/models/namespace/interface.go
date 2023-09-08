package namespace

import (
	"fmt"

	"gorm.io/gorm"
)

// Interface manages the namespace db methods
type Interface interface {
	Create(namespace *Namespace) error
	Delete(namespaceID uint) error
	Update(id uint, namespace *Namespace) error
	GetAll() ([]*Namespace, error)
	GetByID(id uint) (*Namespace, error)
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

func (c core) Update(id uint, namespace *Namespace) error {
	return c.db.Where("id = ?", id).Updates(namespace).Error
}

func (c core) Delete(namespaceID uint) error {
	return c.db.Delete(&Namespace{}, "id = ?", namespaceID).Error
}

func (c core) GetAll() ([]*Namespace, error) {
	list := make([]*Namespace, 0)

	if err := c.db.Find(&list).Error; err != nil {
		return nil, fmt.Errorf("[db.Namespace.Get] failed to get records error=%w", err)
	}

	return list, nil
}

func (c core) GetByID(id uint) (*Namespace, error) {
	namespace := new(Namespace)

	if err := c.db.Preload("Projects").Where("id = ?", id).First(&namespace).Error; err != nil {
		return nil, fmt.Errorf("[db.Namespace.GetByID] failed to get record error=%w", err)
	}

	if namespace.ID != id {
		return nil, ErrRecordNotFound
	}

	return namespace, nil
}
