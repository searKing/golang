package resillience

//go:generate stringer -type Event -trimprefix=Event
//go:generate jsonenums -type Event
type Event int

const (
	EventNew     Event = iota // new and start
	EventClose                // close
	EventExpired              // restart
)
