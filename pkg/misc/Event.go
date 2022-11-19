package misc

type EventType byte

const(
	_ = 				 iota
	EventPut EventType = iota
	EventDelete
)

type Event struct {
	Index 	uint64 	// index for ordering
	Type 	EventType 	// type of event (Put, Delete, etc)
	Key 	string		// key where event happened
	Value 	string	// value associated (with Put)
}