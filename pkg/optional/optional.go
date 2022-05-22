package optional

import "golang.org/x/exp/constraints"

type NonPointer interface {
	constraints.Ordered | ~bool | ~struct{}
}

func Ref[T NonPointer](v T) *T {
	return &v
}
