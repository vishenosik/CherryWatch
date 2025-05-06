package versions

import (
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		major    int
		minor    int
		patch    int
		expected SemanticVersion
	}{
		{"Basic version", 1, 2, 3, SemanticVersion{1, 2, 3}},
		{"Zero version", 0, 0, 0, SemanticVersion{0, 0, 0}},
		{"Large numbers", 999, 888, 777, SemanticVersion{999, 888, 777}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.major, tt.minor, tt.patch)
			if got != tt.expected {
				t.Errorf("New(%d, %d, %d) = %v, want %v", tt.major, tt.minor, tt.patch, got, tt.expected)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		input    string
		expected SemanticVersion
		wantErr  bool
	}{
		{"1.2.3", SemanticVersion{1, 2, 3}, false},
		{"1.5.0", SemanticVersion{1, 5, 0}, false},
		{"0.0.0", SemanticVersion{0, 0, 0}, false},
		{"10.20.30", SemanticVersion{10, 20, 30}, false},
		{"1.2", SemanticVersion{}, true},
		{"1.2.x", SemanticVersion{}, true},
		{"invalid", SemanticVersion{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := tt.expected.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("Parse(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		v1       SemanticVersion
		v2       SemanticVersion
		expected int
	}{
		{New(1, 0, 0), New(1, 0, 0), 0},
		{New(1, 2, 3), New(1, 2, 3), 0},
		{New(1, 0, 0), New(2, 0, 0), -1},
		{New(2, 0, 0), New(1, 0, 0), 1},
		{New(1, 2, 0), New(1, 3, 0), -1},
		{New(1, 3, 0), New(1, 2, 0), 1},
		{New(1, 2, 3), New(1, 2, 4), -1},
		{New(1, 2, 4), New(1, 2, 3), 1},
	}

	for _, tt := range tests {
		t.Run(tt.v1.String()+"_"+tt.v2.String(), func(t *testing.T) {
			got := tt.v1.Compare(tt.v2)
			if got != tt.expected {
				t.Errorf("%v.Compare(%v) = %d, want %d", tt.v1, tt.v2, got, tt.expected)
			}
		})
	}
}

func TestIn(t *testing.T) {
	v := New(1, 2, 3)
	tests := []struct {
		name     string
		v        SemanticVersion
		v1       SemanticVersion
		v2       SemanticVersion
		expected bool
	}{
		{"Within range", v, New(1, 0, 0), New(2, 0, 0), true},
		{"Equal to lower bound", v, v, New(2, 0, 0), true},
		{"Equal to upper bound", v, New(1, 0, 0), v, true},
		{"Below range", v, New(1, 3, 0), New(2, 0, 0), false},
		{"Above range", v, New(0, 9, 0), New(1, 1, 0), false},
		{"Reverse bounds", v, New(2, 0, 0), New(1, 0, 0), true},
		{"Same bounds", v, v, v, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.v.In(tt.v1, tt.v2)
			if got != tt.expected {
				t.Errorf("%v.In(%v, %v) = %v, want %v", tt.v, tt.v1, tt.v2, got, tt.expected)
			}
		})
	}
}

func TestGTE(t *testing.T) {
	tests := []struct {
		v1       SemanticVersion
		v2       SemanticVersion
		expected bool
	}{
		{New(1, 0, 0), New(1, 0, 0), true},
		{New(2, 0, 0), New(1, 0, 0), true},
		{New(1, 2, 0), New(1, 1, 0), true},
		{New(1, 2, 3), New(1, 2, 2), true},
		{New(1, 0, 0), New(2, 0, 0), false},
		{New(1, 1, 0), New(1, 2, 0), false},
		{New(1, 2, 2), New(1, 2, 3), false},
	}

	for _, tt := range tests {
		t.Run(tt.v1.String()+"_"+tt.v2.String(), func(t *testing.T) {
			got := tt.v1.GTE(tt.v2)
			if got != tt.expected {
				t.Errorf("%v.GTE(%v) = %v, want %v", tt.v1, tt.v2, got, tt.expected)
			}
		})
	}
}

func TestLTE(t *testing.T) {
	tests := []struct {
		v1       SemanticVersion
		v2       SemanticVersion
		expected bool
	}{
		{New(1, 0, 0), New(1, 0, 0), true},
		{New(1, 0, 0), New(2, 0, 0), true},
		{New(1, 1, 0), New(1, 2, 0), true},
		{New(1, 2, 2), New(1, 2, 3), true},
		{New(2, 0, 0), New(1, 0, 0), false},
		{New(1, 2, 0), New(1, 1, 0), false},
		{New(1, 2, 3), New(1, 2, 2), false},
	}

	for _, tt := range tests {
		t.Run(tt.v1.String()+"_"+tt.v2.String(), func(t *testing.T) {
			got := tt.v1.LTE(tt.v2)
			if got != tt.expected {
				t.Errorf("%v.LTE(%v) = %v, want %v", tt.v1, tt.v2, got, tt.expected)
			}
		})
	}
}

func TestEQ(t *testing.T) {
	tests := []struct {
		v1       SemanticVersion
		v2       SemanticVersion
		expected bool
	}{
		{New(1, 0, 0), New(1, 0, 0), true},
		{New(1, 2, 3), New(1, 2, 3), true},
		{New(1, 0, 0), New(2, 0, 0), false},
		{New(1, 2, 0), New(1, 3, 0), false},
		{New(1, 2, 3), New(1, 2, 4), false},
	}

	for _, tt := range tests {
		t.Run(tt.v1.String()+"_"+tt.v2.String(), func(t *testing.T) {
			got := tt.v1.EQ(tt.v2)
			if got != tt.expected {
				t.Errorf("%v.EQ(%v) = %v, want %v", tt.v1, tt.v2, got, tt.expected)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		v        SemanticVersion
		expected string
	}{
		{New(1, 2, 3), "1.2.3"},
		{New(0, 0, 0), "0.0.0"},
		{New(10, 20, 30), "10.20.30"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := tt.v.String()
			if got != tt.expected {
				t.Errorf("%v.String() = %q, want %q", tt.v, got, tt.expected)
			}
		})
	}
}
