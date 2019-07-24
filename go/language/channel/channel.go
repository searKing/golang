package channel

func IsClosed(c interface{}) bool {
	vc, ok := c.(chan interface{})
	if !ok {
		return true
	}

	if vc == nil {
		return true
	}

	select {
	case _, ok := <-vc:
		return !ok
	default:
		return false
	}
}
