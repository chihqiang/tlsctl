package webhook

import "testing"

func TestReplacePlaceholders(t *testing.T) {
	cases := []struct {
		name  string
		input string
		vars  map[string]string
		want  string
	}{
		{
			name:  "replace all",
			input: `<__domain__|__cert__|__key__>`,
			vars:  map[string]string{"domain": "example.com", "cert": "CERT", "key": "KEY"},
			want:  `<example.com|CERT|KEY>`,
		},
		{
			name:  "escape json special chars",
			input: `{"c":"__cert__"}`,
			vars:  map[string]string{"cert": "a\"b\nc"},
			want:  `{"c":"a\"b\nc"}`,
		},
		{
			name:  "unknown placeholder kept",
			input: `__nope__`,
			vars:  map[string]string{"domain": "example.com"},
			want:  `__nope__`,
		},
		{
			name:  "no placeholders unchanged",
			input: `plain string`,
			vars:  nil,
			want:  `plain string`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := replacePlaceholders(tc.input, tc.vars); got != tc.want {
				t.Errorf("replacePlaceholders(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
