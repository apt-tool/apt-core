package user

import (
	"fmt"

	"github.com/apt-tool/apt-core/internal/utils/crypto"

	"gorm.io/gorm"
)

// Interface manages the user database methods
type Interface interface {
	Create(user *User) error
	Delete(id uint) error
	Update(id uint, user *User) error
	GetAll() ([]*User, error)
	GetByID(id uint) (*User, error)
	GetByName(name string) (*User, error)
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

func (c core) Create(user *User) error {
	user.Password = crypto.GetMD5Hash(user.Password)

	return c.db.Create(user).Error
}

func (c core) Delete(userID uint) error {
	return c.db.Delete(&User{}, "id = ?", userID).Error
}

func (c core) Update(userID uint, user *User) error {
	return c.db.Where("id = ?", userID).Updates(user).Error
}

func (c core) GetAll() ([]*User, error) {
	list := make([]*User, 0)

	if err := c.db.Find(&list).Error; err != nil {
		return nil, fmt.Errorf("[db.User.Get] failed to get records error=%w", err)
	}

	return list, nil
}

func (c core) GetByID(id uint) (*User, error) {
	user := new(User)

	if err := c.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, fmt.Errorf("[db.User.Get] failed to get records error=%w", err)
	}

	if user.ID != id {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (c core) GetByName(name string) (*User, error) {
	user := new(User)

	if err := c.db.Where("username = ?", name).First(&user).Error; err != nil {
		return nil, fmt.Errorf("[db.User.Get] failed to get records error=%w", err)
	}

	if user.Username != name {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (c core) Validate(name, pass string) (*User, error) {
	user := new(User)

	if err := c.db.Where("username = ?", name).First(&user).Error; err != nil {
		return nil, fmt.Errorf("[db.User.Validate] failed to get user error=%w", err)
	}

	if user.Username != name || user.Password != crypto.GetMD5Hash(pass) {
		return nil, ErrIncorrectPassword
	}

	return user, nil
}
