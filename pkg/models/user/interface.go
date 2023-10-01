package user

import (
	"fmt"

	"github.com/ptaas-tool/base-api/internal/utils/crypto"

	"gorm.io/gorm"
)

// Interface manages the user database methods
type Interface interface {
	Create(user *User) error
	Delete(id uint) error
	GetAll() ([]*User, error)
	Validate(name, pass string) (*User, error)
}

func New(db *gorm.DB) Interface {
	return &core{
		db: db,
	}
}

type core struct {
	db *gorm.DB
}

// Create a new user
func (c core) Create(user *User) error {
	user.Password = crypto.GetMD5Hash(user.Password)

	return c.db.Create(user).Error
}

// Delete an existed user
func (c core) Delete(userID uint) error {
	return c.db.Delete(&User{}, "id = ?", userID).Error
}

// GetAll users
func (c core) GetAll() ([]*User, error) {
	list := make([]*User, 0)

	if err := c.db.Find(&list).Error; err != nil {
		return nil, fmt.Errorf("[db.User.Get] failed to get records error=%w", err)
	}

	return list, nil
}

// Validate a user by its credential
func (c core) Validate(name, pass string) (*User, error) {
	user := new(User)

	if err := c.db.Where("username = ?", name).First(&user).Error; err != nil {
		return nil, ErrUserNotFound
	}

	if user.Username != name || user.Password != crypto.GetMD5Hash(pass) {
		return nil, ErrIncorrectPassword
	}

	return user, nil
}
