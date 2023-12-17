package library

type Library interface {
	// ID returns the ID of the library.
	ID() string
	// List returns the list of media in the library.
	List() ([]Media, error)
}
