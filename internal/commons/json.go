package commons

import "encoding/json"

// UnmarshalGenericJSON unmarshal JSON into given generic type T
func UnmarshalGenericJSON[T any](b []byte) (v T, err error) {
	return v, json.Unmarshal(b, &v)
}
