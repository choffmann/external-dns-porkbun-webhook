package utils

func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0)

	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}

	return result
}

func Map[T, K any](slice []T, fn func(T) K) []K {
	result := make([]K, len(slice))

	for i, item := range slice {
		result[i] = fn(item)
	}

	return result
}

func MapFilter[T, K any](slice []T, fn func(T) *K) []K {
	result := make([]K, 0)

	for _, item := range slice {
    res := fn(item)
    if res != nil {
		  result = append(result, *res)
    }
	}

	return result
}

func Reduce[T, K any](slice []T, initial K, fn func(item T, acc K) K) K {
  current := initial
  for _, e := range slice {
    current = fn(e, current)
  }
  return current
}
