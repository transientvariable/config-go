package config

import "strconv"

// Bool retrieves the boolean value for the provided path.
//
// The returned error will be non-nil if the value corresponding to the provided path:
//   - could not be found
//   - was found, but could not be parsed as a boolean
func Bool(path string) (bool, error) {
	v, err := Value(path)
	if err != nil {
		return false, err
	}
	return strconv.ParseBool(v)
}

// BoolMustResolve is similar behavior to Bool, but panics if an error occurs.
func BoolMustResolve(path string) bool {
	v, err := Bool(path)
	if err != nil {
		panic(err)
	}
	return v
}
