package interfaces

type Finder interface {
	GetQuery() string
	Download() (string, error)
}
