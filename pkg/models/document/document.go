package document

import (
	"github.com/automated-pen-testing/api/pkg/enum"

	"gorm.io/gorm"
)

// Document represents core log files
type Document struct {
	gorm.Model
	ProjectID   uint
	LogFile     string
	Instruction string
	Status      enum.Status
}
