package conversion

import (
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/conversion"
)

func TypeEnforceStruct(t reflect.Type) (reflect.Type, error) {
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type %s is not a struct-type", t)
	}
	return t, nil
}

func TypeEnforceSlice(t reflect.Type) (reflect.Type, error) {
	if t.Kind() != reflect.Slice {
		return nil, fmt.Errorf("type %s is not a slice-type", t)
	}
	return t, nil
}

func CheckStruct(v reflect.Value) error {
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("kind %s is not a struct-kind", v.Kind())
	}
	return nil
}

func EnforceStruct(obj interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(obj)
	if err := CheckStruct(v); err != nil {
		return reflect.Value{}, nil
	}

	return v, nil
}

func extractField(obj interface{}, v reflect.Value, name string) (reflect.Value, error) {
	f := v.FieldByName(name)
	if !f.IsValid() {
		return reflect.Value{}, fmt.Errorf("struct %T has no field %q", obj, name)
	}
	return f, nil
}

func EnforcePtrStruct(obj interface{}) (reflect.Value, error) {
	v, err := conversion.EnforcePtr(obj)
	if err != nil {
		return reflect.Value{}, err
	}

	if err := CheckStruct(v); err != nil {
		return reflect.Value{}, nil
	}
	return v, nil
}

func EnforceStructField(obj interface{}, name string) (reflect.Value, error) {
	v, err := EnforceStruct(obj)
	if err != nil {
		return reflect.Value{}, err
	}

	return extractField(obj, v, name)
}

func EnforcePtrStructField(obj interface{}, name string) (reflect.Value, error) {
	v, err := EnforcePtrStruct(obj)
	if err != nil {
		return reflect.Value{}, err
	}

	return extractField(obj, v, name)
}

func EnforceGetter(obj interface{}, name string) (reflect.Value, error) {
	switch reflect.TypeOf(obj).Kind() {
	case reflect.Struct:
		return EnforceStructField(obj, name)
	case reflect.Ptr:
		return EnforcePtrStructField(obj, name)
	default:
		return reflect.Value{}, fmt.Errorf("cannot obtain getter from non-struct, non-pointer type %T", obj)
	}
}

func EnforceAccessor(obj interface{}, mutable bool, name string) (reflect.Value, error) {
	if !mutable {
		return EnforceGetter(obj, name)
	}
	return EnforcePtrStructField(obj, name)
}

func EnforceSlice(obj interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("object type %T is not a slice-kind", obj)
	}

	return v, nil
}

func EnforceSetSliceIndex(obj interface{}, index int, value interface{}) error {
	v, err := EnforceSlice(obj)
	if err != nil {
		return err
	}

	v.Index(index).Set(reflect.ValueOf(value))
	return nil
}

func EnforceSliceIndex(obj interface{}, index int) (reflect.Value, error) {
	v, err := EnforceSlice(obj)
	if err != nil {
		return reflect.Value{}, err
	}

	return v.Index(index), nil
}

func EnforceSliceIndexGetter(obj interface{}, index int, name string) (reflect.Value, error) {
	v, err := EnforceSliceIndex(obj, index)
	if err != nil {
		return reflect.Value{}, err
	}

	return EnforceGetter(v.Interface(), name)
}

func EnforceSliceIndexAddress(obj interface{}, index int, name string) (reflect.Value, error) {
	v, err := EnforceSlice(obj)
	if err != nil {
		return reflect.Value{}, err
	}

	item := v.Index(index)

	takeAddr := false
	if elemType := v.Type().Elem(); elemType.Kind() != reflect.Ptr && elemType.Kind() != reflect.Interface {
		if !item.CanAddr() {
			return reflect.Value{}, fmt.Errorf("cannot obtain mutable reference to field %q at %d", name, index)
		}
		takeAddr = true
	}

	if takeAddr {
		item = item.Addr().Elem()
	}
	return extractField(item.Interface(), item, name)
}

func EnforceSliceIndexAccessor(obj interface{}, mutable bool, index int, name string) (reflect.Value, error) {
	if !mutable {
		return EnforceSliceIndexGetter(obj, index, name)
	}
	return EnforceSliceIndexAddress(obj, index, name)
}
