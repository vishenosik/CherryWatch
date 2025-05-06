package versions

import (
	"errors"
	"fmt"
)

type SemanticVersion struct {
	Major int
	Minor int
	Patch int
}

// New creates a new SemVersion from major, minor, patch components
func New(major, minor, patch int) SemanticVersion {
	return SemanticVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
	}
}

// String returns the version as "MAJOR.MINOR.PATCH" string
func (v SemanticVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// Parse parses a semantic version string into a SemanticVersion
func (s SemanticVersion) Parse(version string) (SemanticVersion, error) {
	var v SemanticVersion
	_, err := fmt.Sscanf(version, "%d.%d.%d", &v.Major, &v.Minor, &v.Patch)
	if err != nil {
		return SemanticVersion{}, errors.New("version must be format MAJOR.MINOR")
	}
	return v, nil
}

// In checks if the version is between v1 and v2 (inclusive)
func (v SemanticVersion) In(v1, v2 SemanticVersion) bool {
	lower, upper := v1, v2
	if v1.Compare(v2) > 0 {
		lower, upper = upper, lower
	}

	return v.Compare(lower) >= 0 && v.Compare(upper) <= 0
}

// Compare returns:
//
//	-1 if v < other
//	0 if v == other
//	1 if v > other
func (v SemanticVersion) Compare(other SemanticVersion) int {
	if v.Major != other.Major {
		return compareInt(v.Major, other.Major)
	}
	if v.Minor != other.Minor {
		return compareInt(v.Minor, other.Minor)
	}
	return compareInt(v.Patch, other.Patch)
}

// GTE checks if v is greater than or equal to other
func (v SemanticVersion) GTE(other SemanticVersion) bool {
	return v.Compare(other) >= 0
}

// LTE checks if v is less than or equal to other
func (v SemanticVersion) LTE(other SemanticVersion) bool {
	return v.Compare(other) <= 0
}

// EQ checks if v is equal to other
func (v SemanticVersion) EQ(other SemanticVersion) bool {
	return v.Compare(other) == 0
}

// helper function for comparing integers
func compareInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
