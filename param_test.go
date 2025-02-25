package fogg

import "testing"

func TestTagParamHasArg(t *testing.T) {
	param := TagParam{
		Name:  "example",
		Value: "value",
		Args:  []string{"arg1", "arg2", "arg3"},
	}

	tests := []struct {
		arg      string
		expected bool
	}{
		{"arg1", true},
		{"arg2", true},
		{"arg3", true},
		{"arg4", false},
	}

	for _, test := range tests {
		result := param.HasArg(test.arg)
		if result != test.expected {
			t.Errorf("HasArg(%s) = %v; want %v", test.arg, result, test.expected)
		}
	}
}
