package config

import "github.com/multiformats/go-multiaddr"

// Multiaddr retrieves the multiaddr.Multiaddr value for the provided path.
//
// The returned error will be non-nil if the value corresponding to the provided path:
//   - could not be found
//   - was found, but could not be parsed as a multiaddr.Multiaddr value
func Multiaddr(path string) (multiaddr.Multiaddr, error) {
	v, err := Value(path)
	if err != nil {
		return nil, err
	}
	return multiaddr.NewMultiaddr(v)
}

// MultiaddrMustResolve is similar in behavior to Multiaddr, but panics if an error occurs.
func MultiaddrMustResolve(path string) multiaddr.Multiaddr {
	v, err := Multiaddr(path)
	if err != nil {
		panic(err)
	}
	return v
}
