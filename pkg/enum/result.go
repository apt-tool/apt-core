package enum

// Result stands for document result
type Result int

const (
	ResultNotSet Result = iota + 1
	ResultSuccessful
	ResultFailed
	ResultUnknown
)

func (r Result) ToMessage() string {
	switch r {
	case ResultNotSet:
		return "Not set"
	case ResultSuccessful:
		return "Succeed"
	case ResultFailed:
		return "Failed"
	case ResultUnknown:
		return "Unknown"
	}

	return ""
}
