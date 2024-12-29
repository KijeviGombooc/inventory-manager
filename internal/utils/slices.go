package utils

func FirstIndexOf[T any](input []T, targetFunc func(T) bool) int {
	for i, v := range input {
		if targetFunc(v) {
			return i
		}
	}
	return -1
}

func Map[T any, U any](input []T, transform func(T) U) []U {
	result := make([]U, len(input))
	for i, v := range input {
		result[i] = transform(v)
	}
	return result
}

func MapErrored[T any, U any](input []T, transform func(T) (U, error)) ([]U, error) {
	result := make([]U, len(input))
	for i, v := range input {
		transformed, err := transform(v)
		if err != nil {
			return nil, err
		}
		result[i] = transformed
	}
	return result, nil
}

func Reduce[T any, U any](input []T, initial U, reducer func(U, T) U) U {
	accumulator := initial
	for _, v := range input {
		accumulator = reducer(accumulator, v)
	}
	return accumulator
}
