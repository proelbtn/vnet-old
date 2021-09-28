package entities

import "testing"

func TestValidateName(t *testing.T) {
	type arguments string
	type result bool

	tests := []struct {
		name     string
		args     arguments
		expected result
	}{
		{
			name:     "R1",
			args:     "R1",
			expected: true,
		},
		{
			name:     "SW01",
			args:     "SW01",
			expected: true,
		},
		{
			name:     "_",
			args:     "_",
			expected: false,
		},
		{
			name:     "-",
			args:     "-",
			expected: false,
		},
		{
			name:     "123",
			args:     "123",
			expected: true,
		},
		{
			name:     "hoge fuga",
			args:     "hoge fuga",
			expected: false,
		},
		{
			name:     "hoge@fuga",
			args:     "hoge@fuga",
			expected: false,
		},
		{
			name:     "0123456789012345678901234567890123456789",
			args:     "0123456789012345678901234567890123456789",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := validateName(string(tt.args))
			if tt.expected && res != nil {
				t.Errorf("expected res == nil, but returned %v", res)
			}
			if !tt.expected && res == nil {
				t.Errorf("expected res != nil, but returned nil")
			}
		})
	}
}
