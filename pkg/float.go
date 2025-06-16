package config

import "strconv"

// Float retrieves the float value for the provided path.
//
// The returned error will be non-nil and the returned float value will be set to 0 if the value corresponding to the
// provided path:
//   - could not be found
//   - was found, but could not be parsed as a float
func Float(path string) (float64, error) {
	v, err := Value(path)
	if err != nil {
		return 0, err
	}

	if v == "" {
		return 0, nil
	}
	return strconv.ParseFloat(v, 64)
}

// FloatMustResolve is similar behavior to Float, but panics if an error occurs.
func FloatMustResolve(path string) float64 {
	v, err := Float(path)
	if err != nil {
		panic(err)
	}
	return v
}
