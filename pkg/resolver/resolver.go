package resolver

import (
	"reflect"

	"github.com/4strodev/wiring_graphs/pkg/errors"
)

type SimpleFunctionResolver[T any] = func() T
type ErrorFunctionResolver[T any] = func() (T, error)

type DependencyResolver[T any] struct {
	Resolver any
}

func (d DependencyResolver[T]) Type() reflect.Type {
	resolverType := reflect.TypeOf(d.Resolver)
	return resolverType.Out(0)
}

func (d DependencyResolver[T]) Input() []reflect.Type {
	resolverType := reflect.TypeOf(d.Resolver)
	input := []reflect.Type{}
	for i := 0; i < resolverType.NumIn(); i++ {
		input = append(input, resolverType.In(i))
	}

	return input
}

func IsValid(resolver any) bool {
	return canBeSimpleResolver(resolver) || canBeErrorResolver(resolver)
}

func Execute(d DependencyResolver[any], in []reflect.Value) (reflect.Value, error) {
	var value reflect.Value
	var err error

	if canBeSimpleResolver(d.Resolver) {
		reflectValue := reflect.ValueOf(d.Resolver)
		out := reflectValue.Call(in)
		value = out[0]
	} else if canBeErrorResolver(d.Resolver) {
		reflectValue := reflect.ValueOf(d.Resolver)
		out := reflectValue.Call(in)

		if len(out) != 2 {
			err = errors.Errorf(errors.E_INVALID_RESOLVER, "resolver for type %s does not return two values", d.Type().String())
			return reflect.Value{}, err
		}

		value = out[0]
		if out[1].IsZero() {
			err = nil
		}
		returnedError, ok := out[1].Interface().(error)
		if !ok {
			err = errors.Errorf(errors.E_INVALID_RESOLVER, "resolver should return anerror not %s", out[1].Type().String())
		}

		err = returnedError
	}
	return value, err
}

func canBeSimpleResolver(resolver any) bool {
	val := reflect.ValueOf(resolver)
	valType := val.Type()

	if valType.NumOut() != 1 {
		return false
	}

	return true
}

func canBeErrorResolver(resolver any) bool {
	val := reflect.ValueOf(resolver)
	valType := val.Type()

	if valType.NumOut() != 2 {
		return false
	}

	if valType.Out(1) != reflect.TypeFor[error]() {
		return false
	}

	return true
}
