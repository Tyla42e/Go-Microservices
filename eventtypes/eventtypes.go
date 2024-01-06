package eventtypes

type EventType int64

const (
	PostCreated EventType = iota
	CommentCreated
	CommentUpdated
)

func (e EventType) String() string {
	switch e {
	case PostCreated:
		return "PostCreated"
	case CommentCreated:
		return "CommentCreated"
	case CommentUpdated:
		return "CommentUpdated"
	}
	return "Unknown eventtype"
}
