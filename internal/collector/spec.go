package collector

import (
	"reflect"
	"strings"
)

type Spec map[string]Rules

// Rules templates have access to {key}, {old} and {new}.
type Rules struct {
	Changed string `json:"changed,omitzero"`
	Higher  string `json:"higher,omitzero"`
	Lower   string `json:"lower,omitzero"`
}

type FieldSpec struct {
	ptr   any
	rules Rules
}

func Field[V any](ptr *V, rules Rules) FieldSpec {
	return FieldSpec{ptr: ptr, rules: rules}
}

// NewSpec keys each rule by the json name its field marshals with.
func NewSpec(snapshot any, fields ...FieldSpec) Spec {
	v := reflect.ValueOf(snapshot).Elem()
	spec := make(Spec, len(fields))
	for _, f := range fields {
		spec[fieldName(v, f.ptr)] = f.rules
	}
	return spec
}

// fieldName finds which field of the struct ptr points at.
func fieldName(v reflect.Value, ptr any) string {
	for i := range v.NumField() {
		if ptr != v.Field(i).Addr().Interface() {
			continue
		}
		f := v.Type().Field(i)
		if tag, _, _ := strings.Cut(f.Tag.Get("json"), ","); tag != "" {
			return tag
		}
		return f.Name
	}
	panic("collector.NewSpec: rule must target a field of " + v.Type().String())
}
