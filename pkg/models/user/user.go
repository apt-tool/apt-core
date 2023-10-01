package user

import (
	"github.com/ptaas-tool/base-api/pkg/enum"

	"gorm.io/gorm"
)

// User is the base entity of our clients
type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string
	Role     enum.Role
}
