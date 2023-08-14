package user

import (
	"fmt"
	"github.com/automated-pen-testing/api/internal/utils/crypto"
	"gorm.io/gorm"
)

// Interface manages the user database methods
type Interface interface {
	Create(user *User) error
	Delete(userID uint) error
	Update(userID uint, user *User) error
	UpdateInfo(username string, newName string) error
	Get() ([]*User, error)
	GetByID(userID uint) (*User, error)
	GetByIDs(userIDs []uint) ([]*User, error)
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

func (c core) UpdateInfo(username string, newName string) error {
	return c.db.
		Model(&User{}).
		Where("username = ?", username).
		Update("username", newName).
		Error
}

func (c core) Update(userID uint, user *User) error {
	return c.db.Model(&User{}).Update("role", user.Role).Where("id = ?", userID).Error
}

func (c core) Get() ([]*User, error) {
	list := make([]*User, 0)

	if err := c.db.Find(&list).Error; err != nil {
		return nil, fmt.Errorf("[db.User.Get] failed to get records error=%w", err)
	}

	return list, nil
}

func (c core) GetByID(userID uint) (*User, error) {
	user := new(User)

	if err := c.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("[db.User.Get] failed to get records error=%w", err)
	}

	if user.ID != userID {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (c core) GetByIDs(userIDs []uint) ([]*User, error) {
	list := make([]*User, 0)

	if err := c.db.Where("id in ?", userIDs).Find(&list).Error; err != nil {
		return nil, fmt.Errorf("[db.User.Get] failed to get records error=%w", err)
	}

	return list, nil
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
