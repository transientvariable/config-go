package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/transientvariable/anchor"

	"gopkg.in/yaml.v3"
)

const (
	// The default file path to load when creating a new default configuration.
	defaultFilePath = `application.yaml`

	defaultRoot = `config`

	// Format string for matching path prefixes to their corresponding slice elements.
	formatSlicePrefix = `%s\\.#\\d+`

	// Template string for matching placeholder values.
	placeholderTemplate = `$value`

	// Delimiter used for separating configuration placeholder values.
	placeholderValueDelimiter = `|`

	// Regular expression used for matching file extensions.
	fileExtensionPattern = `\.[^.\\/:*?"<>|\r\n]+$`

	// Regular expression used for matching configuration placeholders.
	placeholderPattern = `(?P<value>\${[^}]+})`

	// Regular expression used for matching configuration placeholder values.
	placeholderValuePattern = `\${(?P<value>[^}]+)}`
)

var (
	config  *configuration
	loadErr error
	once    sync.Once
)

type mapConfigFunc func([]byte) (map[string]any, error)

// config is a container for the configuration mapping.
type configuration struct {
	filePath string
	mapping  configMap
	mutex    sync.RWMutex
	root     Path
}

// Load reads and parses the configuration using the provided optional properties.
//
// If an error occurs during read/parse operations, error will be non-nil.
func Load(options ...func(*Option)) error {
	once.Do(func() {
		opts := &Option{}
		for _, opt := range options {
			opt(opts)
		}

		filePath := opts.filePath
		if filePath == "" {
			filePath = defaultFilePath
		}

		rawConfig, err := readConfig(filePath)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				loadErr = err
				return
			}
		}

		mapping, err := newConfigMap(rawConfig)
		if err != nil {
			loadErr = err
			return
		}

		config = &configuration{
			filePath: filePath,
			mapping:  mapping,
			root:     defaultRoot,
		}

		for p, v := range config.mapping {
			config.mapping[p] = interpolate(regexp.MustCompile(placeholderPattern), placeholderTemplate, v)
		}

		r := config.root.String()
		for _, key := range config.mapping.keys() {
			p := strings.Split(key.String(), ".")
			if s := strings.TrimSpace(p[0]); s != "" {
				if !strings.EqualFold(s, r) {
					loadErr = fmt.Errorf("configuration: multiple root paths defined: %s", s)
					break
				}
			}
		}
	})
	return loadErr
}

// hasPath checks whether a configuration value is present for the provided path.
//
// If the path exists, then Value(path) will never result in an error. However, the typed getters, such as
// Int(path), will return a non-nil error if the value is not convertible to the requested type.
func (c *configuration) hasPath(path Path) bool {
	if path.Empty() {
		return false
	}

	path = c.resolve(path)
	for key := range c.mapping {
		if key.Equals(path) {
			return true
		}
	}
	return false
}

func (c *configuration) resolve(path Path) Path {
	if path.Equals(c.root) {
		return path
	}

	if !strings.HasPrefix(path.String(), c.root.String()+".") {
		path = c.root.Join(path)
	}
	return path
}

// isCollection checks whether the configuration value corresponding to the provided Path represents a collection of
// mapping (e.g. slice).
//
// Returns:
//   - true if the value corresponding to the provided Path represents a collection of mapping
//   - false if the value corresponding to the provided Path does not represent a collection of mapping
//   - false if the value corresponding to the Path could not be found
func (c *configuration) isCollection(path Path) bool {
	if !c.hasPath(path) {
		return false
	}
	h := c.hasPath(Path(fmt.Sprintf(formatSliceSuffix, c.resolve(path))))
	return h
}

// set sets or replaces a configuration value for the provided Path.
//
// Returns:
//   - true if the value corresponding to the provided Path was successfully replaced
//   - false if the configuration does not contain the specified Path or otherwise could not replace the value
func (c *configuration) set(path Path, value string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !path.Empty() {
		c.mapping[c.resolve(path)] = value
		return true
	}
	return false
}

// value retrieves the configuration value for the provided path.
//
// The returned error will be non-nil if the value corresponding to the provided path could not be found.
func (c *configuration) value(path Path) (string, error) {
	if !c.hasPath(path) {
		return "", &PathError{Err: ErrPathNotFound, Operation: "value", Path: path.String()}
	}
	return c.mapping[c.resolve(path)], nil
}

// values retrieves the collection of configuration mapping for the provided path.
//
// The returned error will be non-nil if the value corresponding to the path:
//   - could not be found
//   - was found but does not map to a collection of mapping (e.g. slice)
func (c *configuration) values(path Path) ([]string, error) {
	if !c.hasPath(path) {
		return nil, &PathError{Err: ErrPathNotFound, Operation: "values", Path: path.String()}
	}

	if !c.isCollection(path) {
		return nil, &PathError{Err: errors.New("value does not represent a collection"), Path: path.String()}
	}

	var values []string
	slicePrefixPattern := regexp.MustCompile(fmt.Sprintf(formatSlicePrefix, path))
	for p := range c.mapping {
		if slicePrefixPattern.MatchString(p.String()) {
			values = append(values, c.mapping[p])
		}
	}
	return values, nil
}

// Size returns the current number configuration paths.
func Size() int {
	config.mutex.RLock()
	defer config.mutex.RUnlock()
	return len(config.mapping)
}

// String returns a string representation of the configuration.
func (c *configuration) String() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	m := make(map[string]any)
	m["file_path"] = c.filePath
	m["mapping"] = c.mapping
	m["root"] = c.root
	return string(anchor.ToJSONFormatted(m))
}

// HasPath checks whether a configuration value is present for the provided path.
func HasPath(path string) (bool, error) {
	if config == nil {
		return false, fmt.Errorf("configuration: %w", ErrNotInitialized)
	}
	return config.hasPath(Path(path)), nil
}

// Root returns the root configuration Path.
func Root() Path {
	if config == nil {
		panic(fmt.Errorf("configuration: %w", ErrNotInitialized))
	}
	return config.root
}

// Set sets or replaces a configuration value for the provided path.
//
// Returns:
//   - true if the value corresponding to the provided path was successfully replaced
//   - false if the configuration does not contain the specified path or otherwise could not replace the value
func Set(path string, value string) (bool, error) {
	if config == nil {
		return false, fmt.Errorf("configuration: %w", ErrNotInitialized)
	}
	return config.set(Path(path), value), nil
}

// Sub returns the sub-paths for the provided path.
func Sub(path string) ([]Path, error) {
	if config == nil {
		return nil, fmt.Errorf("configuration: %w", ErrNotInitialized)
	}

	hasPath, err := HasPath(path)
	if err != nil {
		return nil, err
	}

	if !hasPath {
		return nil, &PathError{Err: ErrPathNotFound, Operation: "sub", Path: path}
	}

	p := config.resolve(Path(path))
	d := p.Depth() + 1

	var paths []Path
	for key := range config.mapping {
		if strings.Contains(key.String(), p.String()) && key.Depth() == d {
			paths = append(paths, key)
		}
	}
	return paths, nil
}

// String returns a string representation of the configuration.
func String() string {
	if config == nil {
		return fmt.Errorf("configuration: %w", ErrNotInitialized).Error()
	}
	return config.String()
}

func readConfig(filePath string) (map[string]any, error) {
	fileExtension := regexp.MustCompile(fileExtensionPattern).FindString(filePath)
	switch fileExtension {
	case ".json":
		return nil, nil
	case ".yaml", ".yml":
		return readConfigAndThen(filePath, readYaml)
	default:
		return nil, errors.New(fmt.Sprintf(
			"configuration: unsupported file type, expected one of %s, but found %s for path %s",
			[]string{".json", ".yaml", ".yml"}, fileExtension, filePath))
	}
}

func readConfigAndThen(filePath string, mapConfigFn mapConfigFunc) (map[string]any, error) {
	if strings.TrimSpace(filePath) == "" {
		return nil, errors.New("configuration: file path cannot be empty")
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("configuration: %w", err)
	}
	return mapConfigFn(bytes)
}

func readYaml(bytes []byte) (map[string]any, error) {
	var yamlConfig map[string]any
	if err := yaml.Unmarshal(bytes, &yamlConfig); err != nil {
		return nil, fmt.Errorf("configuration: could not read YAML Configuration: %w", err)
	}
	return yamlConfig, nil
}

func interpolate(pattern *regexp.Regexp, template string, value string) string {
	if !pattern.MatchString(value) {
		return value
	}

	for _, match := range findAllMatchesOf(pattern, template, value) {
		placeholderValue := regexp.MustCompile(placeholderValuePattern).FindStringSubmatch(match)[1]
		replacement := interpolateValue(placeholderValue, placeholderValueDelimiter)
		value = strings.Replace(value, match, replacement, -1)
	}
	return strings.TrimSpace(value)
}

func findAllMatchesOf(pattern *regexp.Regexp, template string, value string) []string {
	var matches []byte
	var result []string

	for _, submatches := range pattern.FindAllStringSubmatchIndex(value, -1) {
		submatch := pattern.ExpandString(matches, template, value, submatches)
		result = append(result, string(submatch))
	}
	return result
}

func interpolateValue(value string, delimiter string) string {
	if strings.TrimSpace(value) == "" {
		return value
	}

	replacements := strings.Split(value, delimiter)
	if envReplacement := os.Getenv(strings.TrimSpace(replacements[0])); envReplacement != "" {
		return envReplacement // replace with env variable
	} else if len(replacements) >= 2 {
		return strings.TrimSpace(replacements[1]) // replace with default if provided
	} else {
		return value // replace with value as-is if no suitable replacement is found
	}
}
