package container

import (
	"reflect"
	"strings"

	"github.com/4strodev/wiring/pkg/errors"
)

// This file contains all the logic related to automatically fill structs

func (c Container) Fill(structPointer any) (err error) {
	refStructValue := reflect.ValueOf(structPointer)
	if refStructValue.Kind() != reflect.Pointer {
		return errors.Errorf(errors.E_TYPE_ERROR, "fill expects a struct pointer '%v' was given", refStructValue.Kind())
	}

	if refStructValue.Elem().Kind() != reflect.Struct {
		return errors.Errorf(errors.E_TYPE_ERROR, "fill expects a struct pointer '%v' pointer was given", refStructValue.Elem().Kind())
	}

	refStructType := refStructValue.Elem().Type()

	for i := 0; i < refStructValue.Elem().NumField(); i++ {
		fieldValue := refStructValue.Elem().Field(i)
		if !fieldValue.CanSet() {
			continue
		}

		var instance reflect.Value

		wiringTag := refStructType.Field(i).Tag.Get("wiring")
		if isOmitted(wiringTag) {
			continue
		}
		token := getToken(wiringTag)

		if token != "" {
			instance, err = c.resolveToken(token)
		} else {
			instance, err = c.resolve(fieldValue.Type())
		}
		if err != nil {
			return err
		}

		fieldValue.Set(instance)
	}

	return nil
}

// returns the token that the tag is requesting
// only returns a token when it's found and tag is well
// written. On any other case it just returns an empty string
func getToken(tagValue string) string {
	if tagValue == "" {
		return tagValue
	}

	tagSegments := strings.Split(tagValue, ",")
	// possible values
	// wiring:"tokenName" -> token required
	// wiring:"tokenName,omit" -> token not required
	// wiring:",omit" -> token not required
	if len(tagSegments) > 2 {
		return ""
	}

	if len(tagSegments) == 2 && tagSegments[1] == "omit" {
		return ""
	}

	return strings.Trim(tagSegments[0], " \n\t\r\f")
}

func isOmitted(tagValue string) bool {
	tagSegments := strings.Split(tagValue, ",")
	return len(tagSegments) == 2 && tagSegments[1] == "omit"
}
