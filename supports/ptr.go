package supports

func Val[T any](p *T, fallback ...T) T {
	if p != nil {
		return *p
	}
	if len(fallback) > 0 {
		return fallback[0]
	}
	var zero T
	return zero
}
