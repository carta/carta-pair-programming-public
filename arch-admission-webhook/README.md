# arch-admission-webhook

A **non-blocking** validating admission webhook. On pod `CREATE` it emits one
statsd metric per container image to the cluster Datadog agent, recording
whether the image is multi-arch (amd64 + arm64) and which workload owns it.

Built to feed the arm64 / Graviton migration audit: it surfaces which
workloads still ship single-arch images, as they get admitted.

## Behaviour

- **Always allows** the pod. `failurePolicy: Ignore`, `timeoutSeconds: 3`.
- Registry work is **async** — the admission response returns immediately;
  classification + metric emission happen on a background worker pool. The
  registry can never add latency to pod creation.
- Fail-open: any decode error still returns `Allowed: true`.

## Metric

`k8s.pod.image.arch` (count), tags:

| tag | example |
|-----|---------|
| `multiarch` | `true` / `false` |
| `has_amd64` | `true` / `false` |
| `has_arm64` | `true` / `false` |
| `namespace` | `payments` |
| `controller_kind` | `Deployment` / `Job` / `StatefulSet` / `Pod` |
| `controller_name` | `api-server` |
| `image_repo` | `123.dkr.ecr.us-east-1.amazonaws.com/team/app` |

`image_tag` is deliberately **not** a tag — repo + controller already
identify the workload, and adding the tag would multiply series cardinality
(and Datadog custom-metric cost).

## Current simplifications (by design)

- **Registry lookup is stubbed** (`classify.StubClassifier`): every image is
  reported multi-arch. Real ECR inspection drops in behind the `Classifier`
  interface later (e.g. go-containerregistry); the TTL `Cache` wrapper stays.
- **Owner resolution uses a name-strip heuristic** (`owner.Resolve`): a pod
  owned by a ReplicaSet is attributed to the Deployment whose name is the RS
  name minus its pod-template-hash suffix. No API lookups. Approximate — a
  bare ReplicaSet (no Deployment) is misattributed.

## Coverage caveat

The webhook only sees a workload when a pod is **(re)admitted** — new deploys,
rollouts, scale-ups, restarts. Workloads running untouched since before the
webhook existed are invisible until they restart. For a full audit, pair this
with a periodic scan of controller templates that reuses the same
`classify` package.

## Layout

```
cmd/webhook        server: TLS, /validate, /healthz, async workers
internal/classify  Classifier interface, stub impl, TTL cache
internal/owner     pod ownerRef -> workload (name-strip)
internal/metrics   statsd emitter
internal/webhook   admission handler + async worker pool
deploy/            Deployment, Service, ValidatingWebhookConfiguration
```

## Develop

```
go test ./...
go build ./...
```

## Config (env)

| var | default |
|-----|---------|
| `LISTEN_ADDR` | `:8443` |
| `TLS_CERT_FILE` | `/etc/webhook/certs/tls.crt` |
| `TLS_KEY_FILE` | `/etc/webhook/certs/tls.key` |
| `DD_AGENT_HOST` | `localhost` (set via DownwardAPI `status.hostIP`) |
| `DD_DOGSTATSD_PORT` | `8125` |
