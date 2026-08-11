// Package webhook implements the validating admission handler. It ALWAYS
// allows the pod and never blocks admission — classification and metric
// emission happen asynchronously off the request path.
package webhook

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/carta/arch-admission-webhook/internal/classify"
	"github.com/carta/arch-admission-webhook/internal/metrics"
	"github.com/carta/arch-admission-webhook/internal/owner"
)

type Handler struct {
	Classifier classify.Classifier
	Emitter    *metrics.Emitter
	Log        *slog.Logger
	// work is the async queue; a bounded buffer so a registry stall can
	// never back-pressure the admission path.
	work chan task
}

type task struct {
	pod  corev1.Pod
	ns   string
	ctrl owner.Controller
}

// New starts the async workers and returns the handler.
func New(c classify.Classifier, e *metrics.Emitter, log *slog.Logger, workers, queue int) *Handler {
	h := &Handler{
		Classifier: c,
		Emitter:    e,
		Log:        log,
		work:       make(chan task, queue),
	}
	for i := 0; i < workers; i++ {
		go h.worker()
	}
	return h
}

// ServeHTTP decodes the AdmissionReview, hands the pod to the async queue,
// and immediately returns Allowed. Any decode error still returns Allowed
// (fail-open) to match failurePolicy: Ignore.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeAllow(w, "", "")
		return
	}

	var review admissionv1.AdmissionReview
	if err := json.Unmarshal(body, &review); err != nil || review.Request == nil {
		h.writeAllow(w, "", "")
		return
	}

	var pod corev1.Pod
	uid := review.Request.UID
	if err := json.Unmarshal(review.Request.Object.Raw, &pod); err != nil {
		h.writeAllow(w, uid, review.APIVersion)
		return
	}

	ns := pod.Namespace
	if ns == "" {
		ns = review.Request.Namespace
	}
	ctrl := owner.Resolve(pod.OwnerReferences, pod.Name)

	// Non-blocking enqueue: drop rather than stall if the queue is full.
	select {
	case h.work <- task{pod: pod, ns: ns, ctrl: ctrl}:
	default:
		h.Log.Warn("work queue full, dropping sample", "namespace", ns, "controller", ctrl.Name)
	}

	h.writeAllow(w, uid, review.APIVersion)
}

func (h *Handler) worker() {
	for t := range h.work {
		for _, img := range containerImages(&t.pod) {
			res, err := h.Classifier.Classify(img)
			if err != nil {
				h.Log.Warn("classify failed", "image", img, "err", err)
				continue
			}
			if err := h.Emitter.Emit(metrics.Sample{
				MultiArch:      res.MultiArch,
				HasAMD64:       res.HasAMD64,
				HasARM64:       res.HasARM64,
				Namespace:      t.ns,
				ControllerKind: t.ctrl.Kind,
				ControllerName: t.ctrl.Name,
				ImageRepo:      repoOf(img),
			}); err != nil {
				h.Log.Warn("emit failed", "image", img, "err", err)
			}
		}
	}
}

func containerImages(pod *corev1.Pod) []string {
	var out []string
	for _, c := range pod.Spec.InitContainers {
		out = append(out, c.Image)
	}
	for _, c := range pod.Spec.Containers {
		out = append(out, c.Image)
	}
	return out
}

// repoOf strips the tag and/or digest, leaving registry host + path.
// Handles a registry port (host:port/path) by only treating a ':' after
// the last '/' as a tag separator.
func repoOf(image string) string {
	if i := strings.LastIndex(image, "@"); i >= 0 {
		image = image[:i]
	}
	lastSlash := strings.LastIndex(image, "/")
	if i := strings.LastIndex(image, ":"); i > lastSlash {
		image = image[:i]
	}
	return image
}

func (h *Handler) writeAllow(w http.ResponseWriter, uid types.UID, apiVersion string) {
	resp := admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{Kind: "AdmissionReview", APIVersion: "admission.k8s.io/v1"},
		Response: &admissionv1.AdmissionResponse{UID: uid, Allowed: true},
	}
	if apiVersion != "" {
		resp.APIVersion = apiVersion
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}
