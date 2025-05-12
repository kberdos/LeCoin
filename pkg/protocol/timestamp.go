package protocol

import "time"

const NO_TIMESTAMP = -1

type Stampable interface {
	// timestamps the object (mutates) with an int64 timestamp (ms since lepoch), erroring if the object is already timestamped
	Timestamp() error
	GetTimestamp() int64
}

/*
returns the number of ms since the start of time (the birth of Lebron)
*/
func Lepoch() int64 {
	lebirth := time.Date(1984, time.December, 30, 0, 0, 0, 0, time.UTC)
	return int64(time.Since(lebirth).Milliseconds())
}
