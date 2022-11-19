package log

import "cloud_native_go/pkg/misc"

type TransactionLogger interface {
	LogPut(key, value string)
	LogDelete(key string)
	Err() <-chan error
	ReplayEvents() (<-chan misc.Event, <-chan error)
	Run()
}
