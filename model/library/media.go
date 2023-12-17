package library

type Media interface {
	// ID returns the ID of the media.
	ID() string
	// URI returns the URI of the media.
	URI() string
}
