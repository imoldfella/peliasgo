package gop

func If[T any](t bool, a, b T) T {
	if t {
		return a
	} else {
		return b
	}
}
