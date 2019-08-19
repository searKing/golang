package resilience

type Ptr interface {
	Value() interface{} //actual instance
	Ready() error
	Close()
}
