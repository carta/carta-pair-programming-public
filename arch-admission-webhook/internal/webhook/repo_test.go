package webhook

import "testing"

func TestRepoOf(t *testing.T) {
	tests := []struct {
		image string
		want  string
	}{
		{"nginx", "nginx"},
		{"nginx:1.25", "nginx"},
		{"library/nginx:1.25", "library/nginx"},
		{"123.dkr.ecr.us-east-1.amazonaws.com/team/app:v1", "123.dkr.ecr.us-east-1.amazonaws.com/team/app"},
		{"registry:5000/team/app:v1", "registry:5000/team/app"},
		{"registry:5000/team/app", "registry:5000/team/app"},
		{"team/app@sha256:abcdef", "team/app"},
		{"team/app:v1@sha256:abcdef", "team/app"},
	}
	for _, tt := range tests {
		if got := repoOf(tt.image); got != tt.want {
			t.Errorf("repoOf(%q) = %q, want %q", tt.image, got, tt.want)
		}
	}
}
