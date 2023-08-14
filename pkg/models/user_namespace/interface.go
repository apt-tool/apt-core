package user_namespace

import (
	"fmt"

	"gorm.io/gorm"
)

// Interface handles the relation between users and namespace
type Interface interface {
	Create(namespaceID uint, userIDs []uint) error
	Clear(namespaceID uint) error
	GetUsers(namespaceID uint) ([]uint, error)
	GetNamespaces(userID uint) ([]uint, error)
}

func New(db *gorm.DB) Interface {
	return &core{
		db: db,
	}
}

type core struct {
	db *gorm.DB
}

func (c core) Create(namespaceID uint, userIDs []uint) error {
	for _, id := range userIDs {
		tmp := &UserNamespace{
			NamespaceID: namespaceID,
			UserID:      id,
		}

		if err := c.db.Create(tmp).Error; err != nil {
			return fmt.Errorf("[db.userNamespace.Create] failed to create record error=%w", err)
		}
	}

	return nil
}

func (c core) Clear(namespaceID uint) error {
	return c.db.Delete(&UserNamespace{}, "namespace_id = ?", namespaceID).Error
}

func (c core) GetUsers(namespaceID uint) ([]uint, error) {
	list := make([]uint, 0)

	if err := c.db.Model(&UserNamespace{}).Where("namespace_id = ?", namespaceID).Pluck("user_id", &list).Error; err != nil {
		return nil, fmt.Errorf("[db.userNamespace.Get] failed to get records error=%w", err)
	}

	return list, nil
}

func (c core) GetNamespaces(userID uint) ([]uint, error) {
	list := make([]uint, 0)

	if err := c.db.Model(&UserNamespace{}).Where("user_id = ?", userID).Pluck("namespace_id", &list).Error; err != nil {
		return nil, fmt.Errorf("[db.userNamespace.Get] failed to get records error=%w", err)
	}

	return list, nil
}
