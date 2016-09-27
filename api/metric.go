package api

type Metric struct {
	ReadCount       int64 //lock
	WriteCount      int64 //lock
	ReadBytes       int64 //lock
	ReaderStartTime int64
	ReaderEndTime   int64
	WriterStartTime int64
	WriterEndTime   int64
}
