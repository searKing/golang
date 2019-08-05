package exception

type IllegalStateException struct {
	*RuntimeException
}

func NewIllegalStateException() *IllegalStateException {
	return &IllegalStateException{
		RuntimeException: NewRuntimeException(),
	}
}

func NewIllegalStateException1(message string) *IllegalStateException {
	return &IllegalStateException{
		RuntimeException: NewRuntimeException1(message),
	}
}

func NewIllegalStateException2(message string, cause Throwable) *IllegalStateException {
	return &IllegalStateException{
		RuntimeException: NewRuntimeException2(message, cause),
	}
}
