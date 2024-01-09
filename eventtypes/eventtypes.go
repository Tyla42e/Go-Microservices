package eventtypes

type EventType int64

const (
	PostCreated EventType = iota
	CommentCreated
	CommentUpdated
	CommentModerated
)

func (e EventType) String() string {
	switch e {
	case PostCreated:
		return "PostCreated"
	case CommentCreated:
		return "CommentCreated"
	case CommentUpdated:
		return "CommentUpdated"
	case CommentModerated:
		return "CommentModerated"
	}
	return "Unknown eventtype"
}
