package api

type Record interface {
	Add(v interface{})

	AddWithIndex(k uint, v interface{})

	Get(uint) interface{}

	Size() uint
}
