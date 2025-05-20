package container

import (
	"fmt"
	"reflect"
)

func Resolve[T any](c *Container) (T, error) {
	var dependency T
	_, err := c.DetectCircularDependencies()
	if err != nil {
		return dependency, err
	}

	resolvedValue, err := c.resolve(reflect.TypeFor[T]())
	if err != nil {
		return dependency, err
	}

	var ok bool
	dependency, ok = resolvedValue.Interface().(T)
	if !ok {
		err = fmt.Errorf("cannot convert returned value into %s", reflect.TypeFor[T]().String())
	}

	return dependency, err
}
