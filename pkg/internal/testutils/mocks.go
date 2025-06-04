// this package contains mocks and utilities for the tests across the whole codebase
package testutils

import (
	"bytes"
	"fmt"
)

type MyService struct {
}

func (s MyService) SayHi() string {
	return "hi!"
}

func NewService() MyService {
	return MyService{}
}

type MyDeps struct {
	Service MyService
	Buffer            *bytes.Buffer `wiring:"buffer"`
	autoIgnored       *bytes.Buffer
	ExplicitlyIgnored *bytes.Buffer `wiring:",omit"`
}

func (d MyDeps) CheckResolvedDependencies() error {
	if d.ExplicitlyIgnored != nil {
		return fmt.Errorf("field 'ExplicitlyIgnored' was assigned")
	}

	if d.Buffer == nil {
		return fmt.Errorf("field 'Buffer' was not assigned")
	}

	return nil
}
