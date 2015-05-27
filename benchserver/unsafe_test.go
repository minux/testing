package benchserver

import (
	"reflect"
	"testing"
)

// identicalTypes tests whether two reflect.Types are identical.
// The only exception is that if b comes from the testing package,
// then the pkgpath component of a type doesn't need to match.
func identicalTypes(t *testing.T, a, b reflect.Type) {
	var helper func(a, b reflect.Type)
	visited := make(map[[2]reflect.Type]bool)
	helper = func(a, b reflect.Type) {
		if a == b {
			return
		}
		if visited[[2]reflect.Type{a, b}] {
			return
		}
		visited[[2]reflect.Type{a, b}] = true
		if a.Kind() != b.Kind() {
			t.Fatalf("Kind mismatch for %v and %v", a, b)
		}
		switch kind := a.Kind(); kind {
		default:
			t.Fatalf("unhandled kind %v: %v and %v", kind, a, b)
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
			reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
			reflect.String:
			// ok
		case reflect.Slice, reflect.Ptr:
			helper(a.Elem(), b.Elem())
		case reflect.Map:
			helper(a.Key(), b.Key())
			helper(a.Elem(), b.Elem())
		case reflect.Array:
			if a.Len() != b.Len() {
				t.Fatalf("mismatched array length: %v and %v", a, b)
			}
			helper(a.Elem(), b.Elem())
		case reflect.Chan:
			if a.ChanDir() != b.ChanDir() {
				t.Fatalf("mismatched channel direction: %v and %v", a, b)
			}
			helper(a.Elem(), b.Elem())
		case reflect.Interface:
			if !a.Implements(b) || !b.Implements(a) {
				t.Fatalf("mismatched interface: %v and %v", a, b)
			}
		case reflect.Func:
			if a.IsVariadic() != b.IsVariadic() {
				t.Fatalf("mismatched variadicity: %v and %v", a, b)
			}
			if a.NumIn() != b.NumIn() || a.NumOut() != b.NumOut() {
				t.Fatalf("mismatched parameter count: %v and %v", a, b)
			}
			for i := 0; i < a.NumIn(); i++ {
				helper(a.In(i), b.In(i))
			}
			for i := 0; i < a.NumOut(); i++ {
				helper(a.Out(i), b.Out(i))
			}
		case reflect.Struct:
			if a.NumField() != b.NumField() {
				t.Fatalf("field count mismatch between %v and %v", a, b)
			}
			for i := 0; i < a.NumField(); i++ {
				af, bf := a.Field(i), b.Field(i)
				if af.Name != bf.Name {
					t.Fatalf("field %d name mismatch: %v and %v", i, af, bf)
				}
				if bf.PkgPath != "testing" && af.PkgPath != bf.PkgPath {
					t.Fatalf("field %d pkgPath mismatch: %v and %v", i, af, bf)
				}
				if af.Offset != bf.Offset {
					t.Fatalf("field %d offset mismatch: %v and %v", i, af, bf)
				}
				if !reflect.DeepEqual(af.Index, bf.Index) {
					t.Fatalf("field %d index mismatch: %v and %v", i, af, bf)
				}
				helper(af.Type, bf.Type)
			}
		}
	}
	helper(a, b)
}

func TestB(t *testing.T) {
	identicalTypes(t, reflect.TypeOf(B{}), reflect.TypeOf(testing.B{}))
}

var sink interface{}

func BenchmarkA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sink = new(int)
	}
}
