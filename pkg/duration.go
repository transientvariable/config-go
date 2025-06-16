package config

import "time"

// Duration retrieves the time.Duration value for the provided path.
//
// The returned error will be non-nil if the value corresponding to the provided path:
//   - could not be found
//   - was found, but could not be parsed as a Duration
func Duration(path string) (time.Duration, error) {
	v, err := Value(path)
	if err != nil {
		return 0, err
	}
	return time.ParseDuration(v)
}

// DurationMustResolve is similar behavior to Duration, but panics if an error occurs.
func DurationMustResolve(path string) time.Duration {
	v, err := Duration(path)
	if err != nil {
		panic(err)
	}
	return v
}
