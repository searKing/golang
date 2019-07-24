package exception

type Exception struct {
	*Throwable
}

func NewException() *Exception {
	return &Exception{
		Throwable: NewThrowable(),
	}
}

func NewException1(message string) *Exception {
	return &Exception{
		Throwable: NewThrowable1(message),
	}
}

func NewException2(message string, cause ThrowableInterface) *Exception {
	return &Exception{
		Throwable: NewThrowable2(message, cause),
	}
}

func NewException4(message string, cause *Throwable, enableSuppression, writableStackTrace bool) *Exception {
	return &Exception{
		Throwable: NewThrowable4(message, cause, enableSuppression, writableStackTrace),
	}
}
