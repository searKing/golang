package exception

type NullPointerException struct {
	*Exception
}

func NewNullPointerException(message string) ThrowableInterface {
	return &NullPointerException{}
}
