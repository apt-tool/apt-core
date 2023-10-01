package document

import (
	"time"

	"github.com/ptaas-tool/base-api/pkg/enum"

	"gorm.io/gorm"
)

// Document represents core log files
type Document struct {
	gorm.Model
	ProjectID     uint
	LogFile       string
	Instruction   string
	ExecutedBy    string
	ExecutionTime time.Duration
	Result        enum.Result
	Status        enum.Status
}
