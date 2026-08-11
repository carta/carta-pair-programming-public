// Package metrics emits the per-image arch verdict to the cluster Datadog
// agent over statsd (UDP).
package metrics

import (
	"github.com/DataDog/datadog-go/v5/statsd"
)

const metricName = "k8s.pod.image.arch"

// Emitter sends arch verdicts. image_tag is intentionally NOT a tag: repo +
// controller already identify the workload, and tag would multiply series
// cardinality (and Datadog custom-metric cost).
type Emitter struct {
	client statsd.ClientInterface
}

// New builds an Emitter pointed at addr (e.g. "$DD_AGENT_HOST:8125").
func New(addr string) (*Emitter, error) {
	c, err := statsd.New(addr)
	if err != nil {
		return nil, err
	}
	return &Emitter{client: c}, nil
}

// Sample is one image on one admitted pod.
type Sample struct {
	MultiArch      bool
	HasAMD64       bool
	HasARM64       bool
	Namespace      string
	ControllerKind string
	ControllerName string
	ImageRepo      string
}

func (e *Emitter) Emit(s Sample) error {
	tags := []string{
		"multiarch:" + boolStr(s.MultiArch),
		"has_amd64:" + boolStr(s.HasAMD64),
		"has_arm64:" + boolStr(s.HasARM64),
		"namespace:" + s.Namespace,
		"controller_kind:" + s.ControllerKind,
		"controller_name:" + s.ControllerName,
		"image_repo:" + s.ImageRepo,
	}
	return e.client.Incr(metricName, tags, 1)
}

func (e *Emitter) Close() error { return e.client.Close() }

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
