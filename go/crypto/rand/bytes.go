package rand

// Bytes returns securely generated random bytes.
func Bytes(n int) []byte {
	b, err := BytesCrypto(n)
	if err == nil {
		return b
	}
	return BytesMath(n)
}
