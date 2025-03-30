package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/transientvariable/anchor"
)

// formatSliceSuffix defines the format string for configuration paths that map to elements of a slice.
const formatSliceSuffix = "%s.#"

// configMap represents a map that uses the string type for both keys and values.
type configMap map[Path]string

// newConfigMap creates a new configMap from the provided source.
func newConfigMap(source any) (configMap, error) {
	configMap := make(map[Path]string)
	if source != nil || reflect.ValueOf(source).Kind() == reflect.Map {
		value := reflect.ValueOf(source)
		for _, key := range value.MapKeys() {
			if err := flatten(key.String(), value.MapIndex(key), configMap); err != nil {
				return configMap, err
			}
		}
	}
	return configMap, nil
}

// contains returns true if the map contains the given key.
func (m configMap) contains(key string) bool {
	for _, k := range m.keys() {
		if strings.EqualFold(k.String(), strings.TrimSpace(key)) {
			return true
		}
	}
	return false
}

// delete deletes a key out of the map with the given prefix.
func (m configMap) delete(prefix Path) {
	for k := range m {
		match := k == prefix
		if !match {
			if !strings.HasPrefix(k.String(), prefix.String()) {
				continue
			}

			if k[len(prefix):len(prefix)+1] != "." {
				continue
			}
		}
		delete(m, k)
	}
}

// keys returns all the top-level keys in this map
func (m configMap) keys() []Path {
	keys := make(map[Path]struct{})
	for k := range m {
		idx := strings.Index(k.String(), ".")
		if idx == -1 {
			idx = len(k)
		}
		keys[k[:idx]] = struct{}{}
	}

	result := make([]Path, 0, len(keys))
	for k := range keys {
		result = append(result, k)
	}
	return result
}

// merge merges the contents of the other FlatMapStr into this one.
func (m configMap) merge(source configMap) {
	for _, prefix := range source.keys() {
		m.delete(prefix)
		for k, v := range source {
			if strings.HasPrefix(k.String(), prefix.String()) {
				m[k] = v
			}
		}
	}
}

func flatten(path string, value reflect.Value, data map[Path]string) error {
	if value.Kind() == reflect.Interface {
		value = value.Elem()
	}

	var reflectedValue string
	switch value.Kind() {
	case reflect.Bool:
		if value.Bool() {
			reflectedValue = "true"
		} else {
			reflectedValue = "false"
		}
		break
	case reflect.Int:
		reflectedValue = fmt.Sprintf("%d", value.Int())
		break
	case reflect.Map:
		if err := flattenMap(path, value, data); err != nil {
			return err
		}
		break
	case reflect.Slice:
		if err := flattenSlice(path, value, data); err != nil {
			return err
		}
		break
	case reflect.String:
		reflectedValue = value.String()
		break
	case reflect.Float32, reflect.Float64:
		reflectedValue = value.String()
		break
	default:
		return fmt.Errorf("unknown value type [%s] for path [%s]\nusing data: %s\n", value, path, support.ToJSONFormatted(data))
	}

	data[Path(path)] = reflectedValue
	return nil
}

func flattenMap(path string, value reflect.Value, config map[Path]string) error {
	for _, k := range value.MapKeys() {
		if k.Kind() == reflect.Interface {
			k = k.Elem()
		}

		if k.Kind() != reflect.String {
			panic(fmt.Sprintf("%s: map key is not string: %s", path, k))
		}

		if err := flatten(fmt.Sprintf("%s.%s", path, k.String()), value.MapIndex(k), config); err != nil {
			return err
		}
	}
	return nil
}

func flattenSlice(path string, value reflect.Value, config map[Path]string) error {
	path = fmt.Sprintf(formatSliceSuffix, path)
	config[Path(path)] = fmt.Sprintf("%d", value.Len())
	for i := 0; i < value.Len(); i++ {
		if err := flatten(fmt.Sprintf("%s%d", path, i), value.Index(i), config); err != nil {
			return err
		}
	}
	return nil
}
