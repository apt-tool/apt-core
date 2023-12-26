package enum

type TrackType int

const (
	TrackSuccess TrackType = iota + 1
	TrackWarning
	TrackError
	TrackInProgress
)

func (t TrackType) ToString() string {
	switch t {
	case TrackSuccess:
		return "success"
	case TrackWarning:
		return "warning"
	case TrackError:
		return "danger"
	case TrackInProgress:
		return "primary"
	default:
		return "secondary"
	}
}
