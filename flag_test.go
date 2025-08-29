package envflag

import (
	"bytes"
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
			name:     "testing-123",
			env:      map[string]string{"TESTING_123": "fromenv"},
			expected: "fromenv",
		},
		"flags with env and suffix": {
			name:     "testing-123",
			// Flags take precedence
			values:   []string{"-testing-123", "passed"},
			opts:     []Option{UsageSuffixer()},
			env:      map[string]string{"TESTING_123": "fromenv"},
			expected: "passed",
		},
		"no flags with env and suffix": {
			name:     "testing-123",
			values:   []string{"-h"},
			opts:     []Option{UsageSuffixer()},
			env:      map[string]string{"TESTING_123": "fromenv"},
			expected: "fromenv",
			usage: `Usage of test:
  -testing-123 string
    	usage [TESTING_123]
`,
		},
		"usage prefix": {
			name:     "testing-123",
			values:   []string{"-h"},
			opts:     []Option{UsagePrefixer()},
			env:      map[string]string{"TESTING_123": "fromenv"},
			expected: "fromenv",
			usage: `Usage of test:
  -testing-123 string
    	[TESTING_123] usage
`,
		},
		"usage updater": {
			name:     "testing-123",
			values:   []string{"-h"},
			opts:     []Option{UsageUpdater(func(key,usage string)string{
				return "thine usage"
			})},
			env:      map[string]string{"TESTING_123": "fromenv"},
			expected: "fromenv",
			usage: `Usage of test:
  -testing-123 string
    	thine usage
`,
		},
		"a variety of characters": {
			// Flags take precedence
			name:     "tes_t.ing----123",
			env:      map[string]string{"TES_T_ING_123": "fromenv"},
			expected: "fromenv",
		},
		"prefixed no flag": {
			// Flags take precedence
			name:     "testing-123",
			values:   []string{},
			opts:     []Option{Prefix("APP_")},
			env:      map[string]string{"APP_TESTING_123": "fromenv"},
			expected: "fromenv",
		},
		"prefixed with flag": {
			// Flags take precedence
			name:     "testing-123",
			values:   []string{"-testing-123", "passed"},
			opts:     []Option{Prefix("APP_")},
			env:      map[string]string{"APP_TESTING_123": "fromenv"},
			expected: "passed",
		},
		"no flags custom env": {
			name: "testing-123",
			opts: []Option{Getenv(func(s string) string {
				return "fake env"
			})},
			env:      map[string]string{"TESTING_123": "fromenv"},
			expected: "fake env",
		},
		"no flags with converter": {
			name:     "testing-123",
			opts: []Option{FlagConverter(func(f string) string {
				return "conVERTed"
			})},
			env:      map[string]string{"conVERTed": "fromenv"},
			expected: "fromenv",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			fs := flag.NewFlagSet("test", flag.ContinueOnError)
			buf := new(bytes.Buffer)
			fs.SetOutput(buf)
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
			if tt.usage != buf.String() {
				t.Errorf("got %q, want %q", buf.String(), tt.usage)
			}
		})
	}

}
