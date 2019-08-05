package exception

type IllegalArgumentException struct {
	*RuntimeException
}

func NewIllegalArgumentException() *IllegalArgumentException {
	return &IllegalArgumentException{
		RuntimeException: NewRuntimeException(),
	}
}

func NewIllegalArgumentException1(message string) *IllegalArgumentException {
	return &IllegalArgumentException{
		RuntimeException: NewRuntimeException1(message),
	}
}

func NewIllegalArgumentException2(message string, cause Throwable) *IllegalArgumentException {
	return &IllegalArgumentException{
		RuntimeException: NewRuntimeException2(message, cause),
	}
}
