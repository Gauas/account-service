package supports

import "reflect"

type Patch struct {
	m map[string]interface{}
}

func NewPatch() *Patch {
	return &Patch{m: make(map[string]interface{})}
}

func Set[T any](p *Patch, key string, val *T) *Patch {
	if val != nil {
		p.m[key] = val
	}
	return p
}

func Fill(p *Patch, v any) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		tag := rt.Field(i).Tag.Get("patch")
		if tag == "" || tag == "-" {
			continue
		}
		fv := rv.Field(i)
		if fv.Kind() == reflect.Ptr && !fv.IsNil() {
			p.m[tag] = fv.Interface()
		}
	}
}

func (p *Patch) Build() map[string]interface{} {
	return p.m
}
func (p *Patch) IsEmpty() bool {
	return len(p.m) == 0
}
