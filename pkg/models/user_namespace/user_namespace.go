package user_namespace

import "gorm.io/gorm"

// UserNamespace is the relation between user and a namespace
type UserNamespace struct {
	gorm.Model
	NamespaceID uint
	UserID      uint
}
