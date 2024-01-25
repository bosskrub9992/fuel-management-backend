package validators

import "testing"

type name struct {
	first string `validate:"required"`
	last  string `validate:"required"`
}

func TestRequestValidator_Validate(t *testing.T) {
	n := name{
		first: "",
		last:  "",
	}
	rv := NewRequestValidator()
	if err := rv.Validate(&n); err != nil {
		t.Error(err)
	}
	if err := Validate(&n); err != nil {
		t.Error(err)
	}
}
