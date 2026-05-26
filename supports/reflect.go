package supports

import (
	"errors"
	"reflect"
)

func Fill[T any](dst *T, patch any) error {
	dv, err := structValue(dst, "dst")
	if err != nil {
		return err
	}

	pv, err := structValue(patch, "patch")
	if err != nil {
		return err
	}

	if !pv.IsValid() {
		return nil
	}

	pt := pv.Type()

	for i := 0; i < pt.NumField(); i++ {
		pf := pt.Field(i)

		if !pf.IsExported() {
			continue
		}

		pv := pv.Field(i)

		if skipPatch(pv) {
			continue
		}

		df := dv.FieldByName(pf.Name)

		if skipField(df) {
			continue
		}

		setValue(df, pv.Elem())
	}

	return nil
}

func structValue(v any, name string) (reflect.Value, error) {
	rv := reflect.ValueOf(v)

	if !rv.IsValid() {
		return reflect.Value{}, errors.New(name + " is invalid")
	}

	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return reflect.Value{}, nil
		}

		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return reflect.Value{}, errors.New(name + " must be struct")
	}

	return rv, nil
}

func skipPatch(v reflect.Value) bool {
	return v.Kind() != reflect.Ptr || v.IsNil()
}

func skipField(v reflect.Value) bool {
	return !v.IsValid() || !v.CanSet()
}

func setValue(dst reflect.Value, src reflect.Value) {
	if dst.Kind() == reflect.Ptr {
		setPtr(dst, src)
		return
	}

	setDirect(dst, src)
}

func setPtr(dst reflect.Value, src reflect.Value) {
	t := dst.Type().Elem()

	if src.Type().AssignableTo(t) {
		ptr := reflect.New(t)
		ptr.Elem().Set(src)
		dst.Set(ptr)
		return
	}

	if src.Type().ConvertibleTo(t) {
		ptr := reflect.New(t)
		ptr.Elem().Set(src.Convert(t))
		dst.Set(ptr)
	}
}

func setDirect(dst reflect.Value, src reflect.Value) {
	if src.Type().AssignableTo(dst.Type()) {
		dst.Set(src)
		return
	}

	if src.Type().ConvertibleTo(dst.Type()) {
		dst.Set(src.Convert(dst.Type()))
	}
}
