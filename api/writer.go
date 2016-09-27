package hdata_api

type Writer interface {
	Prepare()
	Execute(Record)
	Close()
}
