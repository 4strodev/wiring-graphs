package container

import (
	"fmt"
	"reflect"
)

// This file contains all the logic related to automatically fill structs

func (c Container) Fill(structPointer any) error {
	refStructValue := reflect.ValueOf(structPointer)
	if refStructValue.Kind() != reflect.Pointer {
		return fmt.Errorf("fill expects a struct pointer '%v' was given", refStructValue.Kind())
	}

	if refStructValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("fill expects a struct pointer '%v' pointer was given", refStructValue.Elem().Kind())
	}

	for i := 0; i < refStructValue.Elem().NumField(); i++ {
		fieldValue := refStructValue.Elem().Field(i)
		if !fieldValue.CanSet() {
			continue
		}
		
		instance, err := c.resolve(fieldValue.Type())
		if err != nil {
			return fmt.Errorf("cannot fill struct: %w", err)
		}

		fieldValue.Set(instance)
	}

	return nil
}
