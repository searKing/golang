package exception

type RuntimeException struct {
	*Exception
}

func NewRuntimeException() *RuntimeException {
	return &RuntimeException{
		Exception: NewException(),
	}
}

func NewRuntimeException1(message string) *RuntimeException {
	return &RuntimeException{
		Exception: NewException1(message),
	}
}

func NewRuntimeException2(message string, cause ThrowableInterface) *RuntimeException {
	return &RuntimeException{
		Exception: NewException2(message, cause),
	}
}

func NewRuntimeException4(message string, cause *Throwable, enableSuppression, writableStackTrace bool) *RuntimeException {
	return &RuntimeException{
		Exception: NewException4(message, cause, enableSuppression, writableStackTrace),
	}
}
