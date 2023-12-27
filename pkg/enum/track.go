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
		return "Successful event"
	case TrackWarning:
		return "Warning"
	case TrackError:
		return "Error event"
	case TrackInProgress:
		return "In-progress event"
	default:
		return "Unknown event"
	}
}
