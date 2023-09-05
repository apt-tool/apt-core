package enum

// Status represents the status of each document
type Status int

const (
	StatusInit Status = iota + 1
	StatusPending
	StatusDone
	StatusFailed
)

func (s Status) ConvertStatusToMessage(status Status) string {
	switch status {
	case StatusInit:
		return "Initialized"
	case StatusPending:
		return "Pending"
	case StatusDone:
		return "Done"
	case StatusFailed:
		return "Failed"
	}

	return ""
}
