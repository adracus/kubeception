package types

import (
	"fmt"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	String         = reflect.TypeOf("")
	MetaV1Duration = reflect.TypeOf(metav1.Duration{})
	MetaV1Time     = reflect.TypeOf(metav1.Time{})
)

type Check interface {
	Check(t reflect.Type) error
}

type CheckFunc func(t reflect.Type) error

func (f CheckFunc) Check(t reflect.Type) error {
	return f(t)
}

func CheckAll(t reflect.Type, checks []Check) error {
	for _, check := range checks {
		if err := check.Check(t); err != nil {
			return err
		}
	}
	return nil
}

func Kind(kind reflect.Kind) Check {
	return CheckFunc(func(t reflect.Type) error {
		if t.Kind() != kind {
			return fmt.Errorf("expected kind of %s to be %s but got %s", t, kind, t.Kind())
		}
		return nil
	})
}

func CheckKind(t reflect.Type, kind reflect.Kind) error {
	if t.Kind() != kind {
		return fmt.Errorf("expected kind of %s to be %s but got %s", t, kind, t.Kind())
	}
	return nil
}

func EqualTo(t reflect.Type) Check {
	return CheckFunc(func(actualT reflect.Type) error {
		if t != actualT {
			return fmt.Errorf("expected type %s but got %s", t, actualT)
		}
		return nil
	})
}
