package api

type Storage interface {
	Put(*Record)

	Puts([]*Record)

	IsEmpty() bool

	Size() int

	Close()
}
