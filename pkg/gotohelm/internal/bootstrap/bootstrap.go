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

const (
	// For reference: https://physics.nist.gov/cuu/Units/binary.html
	milli = 0.001
	kilo  = 1000
	mega  = kilo * kilo
	giga  = kilo * kilo * kilo
	terra = kilo * kilo * kilo * kilo
	peta  = kilo * kilo * kilo * kilo * kilo

	kibi = 1024
	mebi = kibi * kibi
	gibi = kibi * kibi * kibi
	tebi = kibi * kibi * kibi * kibi
	pebi = kibi * kibi * kibi * kibi * kibi
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

func _len(m any) int {
	// Handle empty/nil maps and lists as sprig does not.
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

// pseudo implementation of k8s.io/apimachinery/pkg/api/resource.MustParse.
func resource_MustParse(repr any) any {
	if !TypeIs("string", repr) {
		panic(fmt.Sprintf("invalid Quantity expected string got: %T", repr))
	}

	if !RegexMatch(`^[0-9]+(\.[0-9]){0,1}(k|m|M|G|T|P|Ki|Mi|Gi|Ti|Pi)?$`, repr) {
		panic(fmt.Sprintf("invalid Quantity: %q", repr))
	}

	return repr

	// TODO need to parse and then downcast if we want to have "full" support.
}

func resource_AsInt64(repr any) []any {
	// If repr is numeric, pass it back as it.
	// Probably need to truncate decimals just to be safe...
	if TypeIs("float64", repr) {
		return []any{Int64(repr), true}
	}

	if !TypeIs("string", repr) {
		panic(fmt.Sprintf("resource.Quantity is neither string nor float64: %T", repr))
	}

	// Type cast would work but that relies on bootstrap to work so use sprig
	// functions.
	reprStr := ToString(repr)

	unit := RegexFind("(k|m|M|G|T|P|Ki|Mi|Gi|Ti|Pi)$", reprStr)

	numeric := Float64(Substr(0, len(reprStr)-len(unit), reprStr))

	scale := float64(0)

	if unit == "" {
		scale = 1
	} else if unit == "m" {
		scale = milli
	} else if unit == "k" {
		scale = float64(kilo)
	} else if unit == "M" {
		scale = float64(mega)
	} else if unit == "G" {
		scale = float64(giga)
	} else if unit == "T" {
		scale = float64(terra)
	} else if unit == "P" {
		scale = float64(peta)
	} else if unit == "Ki" {
		scale = float64(kibi)
	} else if unit == "Mi" {
		scale = float64(mebi)
	} else if unit == "Gi" {
		scale = float64(gibi)
	} else if unit == "Ti" {
		scale = float64(tebi)
	} else if unit == "Pi" {
		scale = float64(pebi)
	} else {
		panic(fmt.Sprintf("unknown unit: %q", unit))
	}

	// TODO there's a bug somewhere. Without double casting, we'll get
	// "incompatible types".
	if float64(scale) < float64(1.0) {
		return []any{0, false}
	}

	// TODO There are some cases where AsInt64 returns false. Need to replicate those here...
	// Likely via bounds checks or checks on decimal values.

	// TODO we're possibly losing precision here depending on the unit.
	//
	// if math.MaxFloat64/scale > numeric {
	// 	return []any{0, true}
	// }

	return []any{scale * numeric, true}
}
