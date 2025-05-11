package versions

import (
	"errors"
	"fmt"

	"github.com/vishenosik/CherryWatch/pkg/models"
)

type DoubleVersion struct {
	Major int
	Minor int
}

// Parse parses a semantic version string into a SemanticVersion
func (v DoubleVersion) Parse(version string) (DoubleVersion, error) {
	var double DoubleVersion
	_, err := fmt.Sscanf(version, "%d.%d", &double.Major, &double.Minor)
	if err != nil {
		return DoubleVersion{}, fmt.Errorf("invalid semantic version: %w", err)
	}
	return v, nil
}

// String returns the version as "MAJOR.MINOR" string
func (v DoubleVersion) String() string {
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

func (v DoubleVersion) In(v1, v2 DoubleVersion) bool {
	lower, upper := v1, v2
	if v1.Major > v2.Major || (v1.Major == v2.Major && v1.Minor > v2.Minor) {
		lower, upper = upper, lower
	}

	if lower.Major == upper.Major && lower.Minor == upper.Minor {
		if v.Major != upper.Major && v.Minor != upper.Minor {
			return false
		}
		return true
	} else if (v.Major < lower.Major || v.Major > upper.Major) ||
		(v.Major == lower.Major && v.Minor < lower.Minor) ||
		(v.Major == upper.Major && v.Minor > upper.Minor) {
		return false

	}
	return true
}

func (v DoubleVersion) In_(v1, v2 models.Version) bool {
	converted_v1, ok := v1.(DoubleVersion)
	if !ok {
		return false
	}
	converted_v2, ok := v2.(DoubleVersion)
	if !ok {
		return false
	}

	return v.In(converted_v1, converted_v2)
}

func (v1 DoubleVersion) GTE(v2 DoubleVersion) bool {
	if v1.Major != v2.Major {
		return v1.Major >= v2.Major
	}
	if v1.Minor != v2.Minor {
		return v1.Minor >= v2.Major
	}
	return true
}

// Parse parses a semantic version string into a SemanticVersion
func (v DoubleVersion) Parse_(version string) (models.Version, error) {
	var double DoubleVersion
	_, err := fmt.Sscanf(version, "%d.%d", &double.Major, &double.Minor)
	if err != nil {
		return DoubleVersion{}, errors.New("version must be format MAJOR.MINOR")
	}
	return double, nil
}
