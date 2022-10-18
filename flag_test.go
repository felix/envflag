package envflag

import (
	"flag"
	"os"
	"testing"
)

func TestParseFlagSet(t *testing.T) {
	tests := map[string]struct {
		// Flag name
		name string
		// Default value
		def string
		// Flag command values
		values []string
		opts   []Option
		// The environment
		env map[string]string

		expected string
		usage    string
	}{
		"no flags no env default": {
			name:     "testing-123",
			def:      "defaultvalue",
			expected: "defaultvalue",
		},
		"flags no env no default": {
			name:     "testing-123",
			values:   []string{"-testing-123", "passed"},
			expected: "passed",
		},
		"no flags env": {
			name:     "testing-123",
			env:      map[string]string{"TESTING_123": "fromenv"},
			expected: "fromenv",
		},
		"flags with no env": {
			name:     "testing-123",
			values:   []string{"-testing-123", "passed"},
			expected: "passed",
		},
		"flags with env": {
			// Flags take precedence
			name:     "testing-123",
			values:   []string{"-testing-123", "passed"},
			env:      map[string]string{"TESTING_123": "fromenv"},
			expected: "passed",
		},
		"no flags with env": {
			// Flags take precedence
			name:     "testing-123",
			env:      map[string]string{"TESTING_123": "fromenv"},
			expected: "fromenv",
		},
		"flags with env and suffix": {
			name:     "testing-123",
			values:   []string{"-testing-123", "passed"},
			opts:     []Option{UsageSuffixer()},
			env:      map[string]string{"TESTING_123": "fromenv"},
			expected: "passed",
		},
		"no flags with env and suffix": {
			name:     "testing-123",
			def:      "defaultvalue",
			opts:     []Option{UsageSuffixer()},
			env:      map[string]string{"TESTING_123": "fromenv"},
			expected: "fromenv",
		},
		"a variety of characters": {
			// Flags take precedence
			name:     "tes_t.ing----123",
			env:      map[string]string{"TES_T_ING_123": "fromenv"},
			expected: "fromenv",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ExitOnError)
			for k, v := range tt.env {
				os.Setenv(k, v)
			}
			defer func() {
				for k := range tt.env {
					os.Unsetenv(k)
				}
			}()
			actual := fs.String(tt.name, tt.def, "usage")
			ParseFlagSet(fs, tt.values, tt.opts...)

			if *actual != tt.expected {
				t.Errorf("got %q, want %q", *actual, tt.expected)
			}
		})
	}

}
