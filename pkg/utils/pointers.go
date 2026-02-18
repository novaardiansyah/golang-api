package utils

import "encoding/json"

// String returns a pointer to the string value passed in.
func String(v string) *string {
	return &v
}

// Uint returns a pointer to the uint value passed in.
func Uint(v uint) *uint {
	return &v
}

// Int64 returns a pointer to the int64 value passed in.
func Int64(v int64) *int64 {
	return &v
}

// Bool returns a pointer to the bool value passed in.
func Bool(v bool) *bool {
	return &v
}

// RawMessage returns a pointer to the json.RawMessage value passed in.
func RawMessage(v json.RawMessage) *json.RawMessage {
	return &v
}
