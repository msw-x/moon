package ujson

func NullArray(bytes []byte) []byte {
	if string(bytes) == "null" {
		bytes = []byte("[]")
	}
	return bytes
}
