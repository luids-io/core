// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.md.

// Package option helps to read values from option maps.
//
// This package is a work in progress and makes no API stability promises.
package option

import "fmt"

// Bool returns the field passed as boolean in value, ok if exists.
func Bool(opts map[string]interface{}, field string) (value bool, ok bool, err error) {
	v, ok := opts[field]
	if ok {
		var cast bool
		value, cast = v.(bool)
		if !cast {
			err = fmt.Errorf("invalid '%s'", field)
		}
	}
	return
}

// String returns the field passed as string in value, ok if exists.
func String(opts map[string]interface{}, field string) (value string, ok bool, err error) {
	v, ok := opts[field]
	if ok {
		var cast bool
		value, cast = v.(string)
		if !cast {
			err = fmt.Errorf("invalid '%s'", field)
		}
	}
	return
}

// Int returns the field passed as int in value, ok if exists.
func Int(opts map[string]interface{}, field string) (value int, ok bool, err error) {
	v, ok := opts[field]
	if ok {
		var cast bool
		value, cast = v.(int)
		if cast {
			return
		}
		// when unmarshalling json structs it uses float
		fvalue, cast := v.(float64)
		if !cast {
			err = fmt.Errorf("invalid '%s'", field)
			return
		}
		value = int(fvalue)
	}
	return
}

// Hash returns the field passed as a map of interfaces in value, ok if exists.
func Hash(opts map[string]interface{}, field string) (value map[string]interface{}, ok bool, err error) {
	v, ok := opts[field]
	if ok {
		var cast bool
		value, cast = v.(map[string]interface{})
		if !cast {
			err = fmt.Errorf("invalid '%s'", field)
			return
		}
	}
	return
}

// HashString returns the field passed as a map of strings in value, ok if exists.
func HashString(opts map[string]interface{}, field string) (value map[string]string, ok bool, err error) {
	v, ok := opts[field]
	if ok {
		var cast bool
		value, cast = v.(map[string]string)
		if cast {
			return
		}
		mapiface, cast := v.(map[string]interface{})
		if !cast {
			err = fmt.Errorf("invalid '%s'", field)
			return
		}
		value = make(map[string]string, 0)
		for k, v := range mapiface {
			d, cast := v.(string)
			if !cast {
				err = fmt.Errorf("invalid '%s'", field)
				return
			}
			value[k] = d
		}
	}
	return
}

// SliceString returns the field passed as a slice of strings in value, ok if exists.
func SliceString(opts map[string]interface{}, field string) (value []string, ok bool, err error) {
	v, ok := opts[field]
	if ok {
		var cast bool
		value, cast = v.([]string)
		if cast {
			return
		}
		slice, cast := v.([]interface{})
		if !cast {
			err = fmt.Errorf("invalid '%s'", field)
			return
		}
		value = make([]string, 0, len(slice))
		for _, v := range slice {
			d, cast := v.(string)
			if !cast {
				err = fmt.Errorf("invalid '%s'", field)
				return
			}
			value = append(value, d)
		}
	}
	return
}

// SliceHash returns the field passed as a slice of map of interfaces in value, ok if exists.
func SliceHash(opts map[string]interface{}, field string) (value []map[string]interface{}, ok bool, err error) {
	v, ok := opts[field]
	if ok {
		var cast bool
		value, cast = v.([]map[string]interface{})
		if cast {
			return
		}
		slice, cast := v.([]interface{})
		if !cast {
			err = fmt.Errorf("invalid '%s'", field)
			return
		}
		value = make([]map[string]interface{}, 0, len(slice))
		for _, v := range slice {
			d, cast := v.(map[string]interface{})
			if !cast {
				err = fmt.Errorf("invalid '%s'", field)
				return
			}
			value = append(value, d)
		}
	}
	return
}

// SliceHashString returns the field passed as a slice of map of strings in value, ok if exists.
func SliceHashString(opts map[string]interface{}, field string) (value []map[string]string, ok bool, err error) {
	v, ok := opts[field]
	if ok {
		var cast bool
		value, cast = v.([]map[string]string)
		if cast {
			return
		}

		var slice []map[string]interface{}
		slice, ok, err = SliceHash(opts, field)
		if err != nil {
			return
		}
		if !ok {
			err = fmt.Errorf("invalid '%s'", field)
			return
		}
		value = make([]map[string]string, 0, len(slice))
		for _, vmap := range slice {
			n := make(map[string]string)
			for k, v := range vmap {
				s, cast := v.(string)
				if !cast {
					err = fmt.Errorf("invalid '%s'", field)
					return
				}
				n[k] = s
			}
			value = append(value, n)
		}
	}
	return
}
