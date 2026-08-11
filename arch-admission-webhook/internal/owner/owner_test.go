package owner

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ptr(b bool) *bool { return &b }

func TestResolve(t *testing.T) {
	tests := []struct {
		name    string
		owners  []metav1.OwnerReference
		podName string
		want    Controller
	}{
		{
			name:    "no owner falls back to pod",
			owners:  nil,
			podName: "loose-pod",
			want:    Controller{Kind: "Pod", Name: "loose-pod"},
		},
		{
			name:   "replicaset strips hash to deployment",
			owners: []metav1.OwnerReference{{Kind: "ReplicaSet", Name: "api-server-5f6b8c9d4", Controller: ptr(true)}},
			want:   Controller{Kind: "Deployment", Name: "api-server"},
		},
		{
			name:   "deployment name with dashes keeps all but hash",
			owners: []metav1.OwnerReference{{Kind: "ReplicaSet", Name: "my-api-server-7d4bcf9", Controller: ptr(true)}},
			want:   Controller{Kind: "Deployment", Name: "my-api-server"},
		},
		{
			name:   "job owns pod directly",
			owners: []metav1.OwnerReference{{Kind: "Job", Name: "nightly-batch", Controller: ptr(true)}},
			want:   Controller{Kind: "Job", Name: "nightly-batch"},
		},
		{
			name: "picks the controller owner among many",
			owners: []metav1.OwnerReference{
				{Kind: "SomethingElse", Name: "sidecar"},
				{Kind: "StatefulSet", Name: "kafka", Controller: ptr(true)},
			},
			want: Controller{Kind: "StatefulSet", Name: "kafka"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Resolve(tt.owners, tt.podName)
			if got != tt.want {
				t.Errorf("Resolve() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
