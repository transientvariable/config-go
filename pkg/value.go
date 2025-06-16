package config

import "fmt"

// Value retrieves the value for the provided path.
//
// The returned error will be non-nil if:
//   - the configuration has not been initialized
//   - the value corresponding to the provided path could not be found
func Value(path string) (string, error) {
	if config == nil {
		return "", fmt.Errorf("configuration: %w", ErrNotInitialized)
	}
	return config.value(Path(path))
}

// ValueMustResolve is similar behavior to Value, but panics if an error occurs.
func ValueMustResolve(path string) string {
	v, err := Value(path)
	if err != nil {
		panic(err)
	}
	return v
}

// ValuesMustResolve is similar behavior to values, but panics if an error occurs.
func ValuesMustResolve(path string) []string {
	v, err := config.values(Path(path))
	if err != nil {
		panic(err)
	}
	return v
}