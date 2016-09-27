package api

type Reader interface {
	Prepare()

	Execute()

	Close()
}
