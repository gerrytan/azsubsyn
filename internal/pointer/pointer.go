package pointer

func From[T any](input *T) T {
	var v T
	if input != nil {
		return *input
	}
	return v
}

func To[T any](input T) *T {
	return &input
}
