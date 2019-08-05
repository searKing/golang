package exception

type IndexOutOfBoundsException struct {
	*RuntimeException
}

func NewIndexOutOfBoundsException() *IndexOutOfBoundsException {
	return &IndexOutOfBoundsException{
		RuntimeException: NewRuntimeException(),
	}
}

func NewIndexOutOfBoundsException1(message string) *IndexOutOfBoundsException {
	return &IndexOutOfBoundsException{
		RuntimeException: NewRuntimeException1(message),
	}
}

func NewIndexOutOfBoundsException2(message string, cause Throwable) *IndexOutOfBoundsException {
	return &IndexOutOfBoundsException{
		RuntimeException: NewRuntimeException2(message, cause),
	}
}
