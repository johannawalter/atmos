package utils

import "github.com/samber/lo"

// Coalesce returns the first non-empty arguments. Arguments must be comparable
func Coalesce[T comparable](v ...T) (result T) {
	result, _ = lo.Coalesce(v...)
	return
}
