package config

import "net/url"

// URL retrieves the url.URL value for the provided path.
//
// The returned error will be non-nil if the value corresponding to the provided path:
//   - could not be found
//   - was found, but could not be parsed as a url.URL value
func URL(path string) (*url.URL, error) {
	v, err := Value(path)
	if err != nil {
		return nil, err
	}
	return url.Parse(v)
}

// URLMustResolve is similar in behavior to URL, but panics if an error occurs.
func URLMustResolve(path string) *url.URL {
	v, err := URL(path)
	if err != nil {
		panic(err)
	}
	return v
}
