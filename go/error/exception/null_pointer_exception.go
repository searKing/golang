package exception

type NullPointerException struct {
	*Exception
}

func NewNullPointerException() Throwable {
	return &NullPointerException{
		Exception: NewException(),
	}
}

func NewNullPointerException1(message string) Throwable {
	return &NullPointerException{
		Exception: NewException1(message),
	}
}
