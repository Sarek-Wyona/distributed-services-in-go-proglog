package log

//Config is needed to configure the max size of a segment’s store and index. It helps to centralize the log’s
//configuration, making it easy to configure the log and use the configs throughout the code
type Config struct {
	Segment struct {
		MaxStoreBytes uint64
		MaxIndexBytes uint64
		InitialOffset uint64
	}
}
