package container

import (
	"reflect"

	"github.com/4strodev/wiring_graphs/pkg/errors"
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
		err = errors.Errorf(errors.E_TYPE_ERROR, "cannot convert returned value into %s", reflect.TypeFor[T]().String())
	}

	return dependency, err
}

func ResolveToken[T any](c *Container, token string) (T, error) {
	var dependency T
	_, err := c.DetectCircularDependencies()
	if err != nil {
		return dependency, err
	}

	resolvedValue, err := c.resolveToken(token)
	if err != nil {
		return dependency, err
	}

	var ok bool
	dependency, ok = resolvedValue.Interface().(T)
	if !ok {
		err = errors.Errorf(errors.E_TYPE_ERROR, "cannot convert returned value into %s", reflect.TypeFor[T]().String())
	}

	return dependency, err
}
