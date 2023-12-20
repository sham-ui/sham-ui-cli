package graceful_shutdown

type Priority uint8

const (
	PriorityHigh Priority = iota + 1
	PriorityMedium
)
