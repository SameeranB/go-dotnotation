package dotnotation

import "fmt"

type PropertyNotFoundInMapError struct {
	property string
}

func (e PropertyNotFoundInMapError) Error() string {
	return fmt.Sprintf("cannot get non-existent property '%s' on a map", e.property)
}
