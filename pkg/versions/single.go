package versions

import "fmt"

type SingleVersion int64

// String returns the version as "MAJOR.MINOR" string
func (v SingleVersion) String() string {
	return fmt.Sprintf("%d", v)
}

func (v SingleVersion) In(v1, v2 SingleVersion) bool {
	if v1 > v2 {
		return v2 <= v && v <= v1
	}
	return v1 <= v && v <= v2
}

func (v SingleVersion) GTE(v1 SingleVersion) bool {
	return v >= v1
}
