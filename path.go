package config

import (
	"encoding"
	"strings"
)

var _ encoding.TextMarshaler = (*Path)(nil)

// Path represents a Configuration path.
type Path string

// Base returns last path element for the Configuration Path. For example, if Path.String returns `foo.bar.gopher`,
// then Base would return `gopher`.
func (p Path) Base() string {
	e := strings.Split(p.String(), ".")
	if len(e) == 1 {
		return e[0]
	}
	return e[len(e)-1]
}

// Depth returns the number of elements contained in the Path.
func (p Path) Depth() int {
	return len(strings.Split(p.String(), "."))
}

// Empty returns whether the Path is empty
func (p Path) Empty() bool {
	return p.String() == ""
}

// Equals returns whether the provided path equals Path.
func (p Path) Equals(path Path) bool {
	return strings.EqualFold(p.String(), path.String())
}

// Join joins Path with the provided Path.
func (p Path) Join(path Path) Path {
	if p.Empty() {
		return path
	}

	if path.Empty() {
		return p
	}
	return Path(p.String() + "." + path.String())
}

// MarshalText implements encoding.TextMarshaler for the Path.
func (p Path) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

// String returns the raw string value for the Configuration Path.
func (p Path) String() string {
	return strings.TrimSpace(string(p))
}
