// Welcome to the magical bootstrap package. This package/file generates the
// _shims.tpl file included in all gotohelm outputs. A task in Taskfile.yaml is
// used to copy the generated file into the gotohelm package. In the future, it
// might be easier to transpile this file on the fly.
//
// Because this file sets up basic utilities and bridges between go and
// templating there are restricts on what may be used.
//
//   - go primitives without direct template support (switches, multi-value
//     returns, type assertions, etc) may not be used.
//   - Only go builtins with direct template support (fmt.Sprintf, etc) may be
//     called/imported.
//   - sprig functions must have a binding declared in sprig.go.
//
// +gotohelm:filename=_shims.tpl
// +gotohelm:namespace=_shims
package bootstrap

import (
	"fmt"
)

// isIntLikeFloat is a workaround for JSON always representing numbers as
// float64's. If a value is a float64 with no fractional value, it's considered
// to be an "integer like" float and therefore will pass when type checked via
// typetest or typeassertion.
func isIntLikeFloat(value any) bool {
	// Could also try doing something funky with Printf?
	return TypeIs("float64", value) && (Float64(value)-Floor(value)) == float64(0)
}

// typeatest is the implementation of the go syntax `_, _ := m.(t)`.
func typetest(typ string, value, zero any) []any {
	if TypeIs(typ, value) {
		return []any{value, true}
	}
	return []any{zero, false}
}

// typeassertion is the implementation of the go syntax `_ := m.(t)`.
func typeassertion(typ string, value any) any {
	// canCastToInt := isIntLikeFloat(value)
	canCastToInt := TypeIs("float64", value) && (Float64(value)-Floor(value)) == float64(0)

	if typ == "int" && canCastToInt {
		return Int(value)
	} else if typ == "int32" && canCastToInt {
		return Int(value)
	} else if typ == "int64" && canCastToInt {
		return Int64(value)
	}

	if !TypeIs(typ, value) {
		panic(fmt.Sprintf("expected type of %q got: %T", typ, value))
	}
	return value
}

// dicttest is the implementation of the go syntax `_, _ := m[k]`.
func dicttest(m map[string]any, key string, zero any) []any {
	if HasKey(m, key) {
		return []any{m[key], true}
	}
	return []any{zero, false}
}

// compact is the implementation of `helmette.CompactN`.
// It's a strange and hacky way of handling multi-value returns.
func compact(args []any) map[string]any {
	out := map[string]any{}
	for i, e := range args {
		out[fmt.Sprintf("T%d", 1+i)] = e
	}
	return out
}

// deref is the implementation of the go syntax `*variable`.
func deref(ptr any) any {
	if ptr == nil {
		panic("nil dereference")
	}
	return ptr
}

func len(m map[string]any) int {
	if m == nil {
		return 0
	}
	return Len(m)
}

// re-implementation of k8s.io/utils/ptr.Deref.
func ptr_Deref(ptr, def any) any {
	if ptr != nil {
		return ptr
	}
	return def
}

// re-implementation of k8s.io/utils/ptr.Equal.
func ptr_Equal(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	return a == b
}
