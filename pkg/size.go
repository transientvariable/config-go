package config

import "github.com/dustin/go-humanize"

// SizeBytes retrieves the value representing a unit value of bytes for the provided path.
//
// The returned error will be non-nil if the value corresponding to the provided path:
//   - could not be found
//   - was found, but could not be parsed as a byte size value
func SizeBytes(path string) (int64, error) {
	v, err := Value(path)
	if err != nil {
		return 0, err
	}

	if v == "" {
		return 0, nil
	}

	s, err := humanize.ParseBytes(v)
	if err != nil {
		return 0, err
	}
	return int64(s), nil
}

// SizeBytesMustResolve is similar behavior to SizeBytes, but panics if an error occurs.
func SizeBytesMustResolve(path string) int64 {
	v, err := SizeBytes(path)
	if err != nil {
		panic(err)
	}
	return v
}
