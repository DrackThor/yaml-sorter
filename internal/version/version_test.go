package ysortversion

import "testing"

func TestFormatLocalVersion(t *testing.T) {
	tests := []struct {
		name   string
		tag    string
		commit string
		want   string
	}{
		{
			name:   "tag and commit",
			tag:    "v1.2.3",
			commit: "abc1234",
			want:   "v1.2.3-local.abc1234",
		},
		{
			name:   "missing tag falls back",
			tag:    "",
			commit: "abc1234",
			want:   "v0.0.0-local.abc1234",
		},
		{
			name:   "missing commit falls back",
			tag:    "v1.2.3",
			commit: "",
			want:   "v1.2.3-local.unknown",
		},
		{
			name:   "trims whitespace",
			tag:    " v1.2.3 ",
			commit: " abc1234 ",
			want:   "v1.2.3-local.abc1234",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := formatLocalVersion(tc.tag, tc.commit)
			if got != tc.want {
				t.Fatalf("formatLocalVersion(%q, %q) = %q, want %q", tc.tag, tc.commit, got, tc.want)
			}
		})
	}
}
