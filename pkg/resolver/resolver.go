package resolver

import (
	"reflect"
)

type SimpleFunctionResolver[T any] = func() T
type ErrorFunctionResolver[T any] = func() (T, error)

type DependencyBuilder[T any] struct {
	Resolver any
}

func Execute[T any](d DependencyBuilder[T]) (T, error) {
	var value T
	var err error

	if simpleResolver, ok := isSimpleResolver[T](d.Resolver); ok {
		value = simpleResolver()
	} else if errorResolver, ok := isErrorResolver[T](d.Resolver); ok {
		value, err = errorResolver()
	}

	return value, err
}

func isSimpleResolver[T any](resolver any) (simpleResolver SimpleFunctionResolver[T], isValid bool) {
	val := reflect.ValueOf(resolver)
	valType := val.Type()

	if valType.NumOut() != 1 {
		return
	}

	if valType.Out(0) != reflect.TypeFor[T]() {
		return
	}

	isValid = true
	simpleResolver = func() T {
		var resolverValue T
		result := val.Call([]reflect.Value{})
		value := result[0]

		resolverValue = value.Interface().(T)

		return resolverValue
	}

	return
}

func isErrorResolver[T any](resolver any) (errorResolver ErrorFunctionResolver[T], isValid bool) {
	val := reflect.ValueOf(resolver)
	valType := val.Type()

	if valType.NumOut() != 2 {
		return
	}

	if valType.Out(0) != reflect.TypeFor[T]() {
		return
	}

	if valType.Out(1) != reflect.TypeFor[error]() {
		return
	}

	isValid = true
	errorResolver = func() (T, error) {
		var resolverValue T
		result := val.Call([]reflect.Value{})
		value := result[0]
		err := result[1]

		if !err.IsNil() {
			return resolverValue, err.Interface().(error)
		}

		resolverValue = value.Interface().(T)

		return resolverValue, nil
	}

	return
}
