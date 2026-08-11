// Package owner resolves the workload controller that owns a pod.
//
// It uses a name-strip heuristic: a pod owned by a ReplicaSet is attributed
// to the Deployment whose name is the ReplicaSet name minus its pod-template
// hash suffix. No API lookups. This is cheap but approximate — a bare
// ReplicaSet (no Deployment) is misattributed, and hand-crafted names that
// look like a hash suffix can be over-stripped.
package owner

import (
	"regexp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Controller is the workload a pod is attributed to.
type Controller struct {
	Kind string
	Name string
}

// ReplicaSet names are "<deployment>-<pod-template-hash>". The hash is an
// alphanumeric segment, commonly 8-10 chars. Strip exactly one trailing
// "-<hash>" segment.
var rsHashSuffix = regexp.MustCompile(`-[a-z0-9]+$`)

// Resolve picks the controlling ownerReference and maps it to a workload.
// Falls back to kind=Pod when there is no controller ownerReference.
func Resolve(owners []metav1.OwnerReference, podName string) Controller {
	ref := controllerRef(owners)
	if ref == nil {
		return Controller{Kind: "Pod", Name: podName}
	}

	if ref.Kind == "ReplicaSet" {
		return Controller{Kind: "Deployment", Name: rsHashSuffix.ReplaceAllString(ref.Name, "")}
	}

	// Job, StatefulSet, DaemonSet, etc. own pods directly.
	return Controller{Kind: ref.Kind, Name: ref.Name}
}

func controllerRef(owners []metav1.OwnerReference) *metav1.OwnerReference {
	for i := range owners {
		if owners[i].Controller != nil && *owners[i].Controller {
			return &owners[i]
		}
	}
	if len(owners) > 0 {
		return &owners[0]
	}
	return nil
}
