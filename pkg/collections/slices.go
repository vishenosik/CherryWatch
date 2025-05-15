package collections

import (
	"iter"
)

func Iter[S ~[]T, T any](slice S) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, elem := range slice {
			if !yield(elem) {
				return
			}
		}
	}
}

func FilterCount[T any](seq iter.Seq[T], by func(T) bool) (it iter.Seq[T], cnt int) {
	it = Filter(seq, by)
	for range it {
		cnt++
	}
	return it, cnt
}

func Filter[T any](seq iter.Seq[T], by func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := range seq {
			if !by(i) {
				return
			}
			if !yield(i) {
				return
			}
		}
	}
}

func Unique[Slice ~[]Type, Type comparable](slice Slice) Slice {
	buffer := make(map[Type]struct{}, len(slice))
	for _, elem := range slice {
		buffer[elem] = struct{}{}
	}
	unique := make(Slice, 0, len(buffer))
	for elem := range buffer {
		unique = append(unique, elem)
	}
	return unique
}

func HasDuplicates[Type comparable](slice ...Type) bool {
	return len(slice) != len(Unique(slice))
}

func ConvertSlice[Type1 any, Type2 any](t1 []Type1, converter func(Type1) Type2) []Type2 {
	out := make([]Type2, 0, len(t1))
	for _, t := range t1 {
		out = append(out, converter(t))
	}
	return out
}
