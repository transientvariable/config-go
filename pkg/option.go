package config

import "strings"

// Option is a container for optional properties that can be used for initializing the configuration.
type Option struct {
	filePath string
}

// WithFilePath sets the file path Option for the configuration. If the file path is not provided, the root of the
// project directory will be used for loading the configuration.
func WithFilePath(filePath string) func(*Option) {
	return func(o *Option) {
		o.filePath = strings.TrimSpace(filePath)
	}
}
