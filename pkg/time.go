package config

import (
	"time"

	"github.com/timberio/go-datemath"
)

// Time retrieves the time.Time value for the provided path.
//
// The returned error will be non-nil if the value corresponding to the provided path:
//   - could not be found
//   - was found, but could not be parsed as a time.Time value
func Time(path string) (time.Time, error) {
	v, err := Value(path)
	if err != nil {
		return time.Time{}, err
	}

	expr, err := datemath.Parse(v)
	if err != nil {
		return time.Time{}, err
	}
	return expr.Time(), nil
}

// TimeMustResolve is similar behavior to Time, but panics if an error occurs.
func TimeMustResolve(path string) time.Time {
	v, err := Time(path)
	if err != nil {
		panic(err)
	}
	return v
}
