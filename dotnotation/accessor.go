package dotnotation

import (
	"errors"
	"fmt"
)

// Accessor provides two methods, Get and Set, that can be configured to handle custom data structures via the
// exported properties, Parser, Getter, and Setter.
type Accessor struct {
	// Getter returns the property value of a given target, or an error.
	Getter func(target interface{}, property string) (interface{}, error)
	// Setter sets the property value of a given target, to a given value, or returns an error.
	Setter func(target interface{}, property string, value interface{}) error
	// Parser converts a given key into a list of properties to access in order to get or set.
	Parser func(key string) []string
	// ValueParser type casts the value received from a Get call.
	ValueParser func(value interface{}) (interface{}, error)
}

func (p Accessor) Set(target interface{}, key string, value interface{}) error {
	properties := p.parser(key)

	for i, property := range properties {
		if i == (len(properties) - 1) {
			// we reached the last property
			return p.setter(target, property, value)
		}

		// create the missing property if it does not exist
		if _, ok := target.(map[string]interface{})[property]; !ok {
			if m, ok := target.(map[string]interface{}); !ok {
				return fmt.Errorf("type conversion failed")
			} else {
				m[property] = map[string]interface{}{}
				target = m[property]
			}
		}
	}
	return errors.New("no properties parsed from key: " + key)
}

func (p Accessor) Get(target interface{}, key string) (interface{}, error) {
	properties := p.parser(key)

	for i, property := range properties {
		if i == (len(properties) - 1) {
			// we reached the last property
			return p.getter(target, property)
		}

		// attempt to get the next level
		var err error
		target, err = p.getter(target, property)
		if err != nil {
			return nil, err
		}
	}

	return nil, errors.New("no properties parsed from key: " + key)
}

func (p Accessor) getter(target interface{}, property string) (interface{}, error) {
	if p.Getter == nil {
		val, err := DefaultGetter(target, property)
		if err != nil {
			return nil, err
		}
		return p.valueParser(val)
	}

	return p.Getter(target, property)
}

func (p Accessor) setter(target interface{}, property string, value interface{}) error {
	if p.Setter == nil {
		return DefaultSetter(target, property, value)
	}

	return p.Setter(target, property, value)
}

func (p Accessor) parser(key string) []string {
	if p.Parser == nil {
		return DefaultParser(key)
	}

	return p.Parser(key)
}

func (p Accessor) valueParser(value interface{}) (interface{}, error) {
	if p.ValueParser == nil {
		return DefaultValueParser(value)
	}

	return p.ValueParser(value)
}
