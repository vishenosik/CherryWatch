// TODO: move to web-tools
package collections

// ToDo: use github.com/vishenosik/web
func ConvertSlice[Type1 any, Type2 any](t1 []Type1, converter func(Type1) Type2) []Type2 {
	out := make([]Type2, 0, len(t1))
	for _, t := range t1 {
		out = append(out, converter(t))
	}
	return out
}
