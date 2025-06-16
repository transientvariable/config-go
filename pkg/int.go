package config

import "strconv"

// Int retrieves the integer value for the provided path.
//
// The returned error will be non-nil and the returned integer value will be set to 0 if the value corresponding to the
// provided path:
//   - could not be found
//   - was found, but could not be parsed as an integer
func Int(path string) (int, error) {
	v, err := Value(path)
	if err != nil {
		return 0, err
	}

	if v == "" {
		return 0, nil
	}
	return strconv.Atoi(v)
}

// IntMustResolve is similar behavior to Int, but panics if an error occurs.
func IntMustResolve(path string) int {
	v, err := Int(path)
	if err != nil {
		panic(err)
	}
	return v
}
