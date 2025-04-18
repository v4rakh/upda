package json

import enc "encoding/json"

// UnmarshalGenericJSON unmarshal JSON into given generic type T
func UnmarshalGenericJSON[T any](b []byte) (v T, err error) {
	return v, enc.Unmarshal(b, &v)
}
