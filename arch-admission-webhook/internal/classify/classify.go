// Package classify determines whether a container image is multi-arch
// (has both linux/amd64 and linux/arm64 in its manifest index).
//
// Registry inspection is intentionally stubbed for now: every image is
// reported as multi-arch. The real ECR lookup drops in behind the
// Classifier interface later; the cache wrapper stays as-is.
package classify

import (
	"sync"
	"time"
)

// Result is the arch verdict for a single image reference.
type Result struct {
	MultiArch bool
	HasAMD64  bool
	HasARM64  bool
}

// Classifier resolves an image reference to an arch verdict.
type Classifier interface {
	Classify(imageRef string) (Result, error)
}

// StubClassifier assumes every image is multi-arch. Placeholder until the
// go-containerregistry ECR lookup is wired in.
type StubClassifier struct{}

func (StubClassifier) Classify(imageRef string) (Result, error) {
	return Result{MultiArch: true, HasAMD64: true, HasARM64: true}, nil
}

// Cache wraps a Classifier with a TTL map keyed by image reference.
// Many pods share an image, so this collapses registry work once the
// real lookup exists.
type Cache struct {
	inner Classifier
	ttl   time.Duration

	mu      sync.Mutex
	entries map[string]entry
}

type entry struct {
	result  Result
	expires time.Time
}

func NewCache(inner Classifier, ttl time.Duration) *Cache {
	return &Cache{
		inner:   inner,
		ttl:     ttl,
		entries: make(map[string]entry),
	}
}

func (c *Cache) Classify(imageRef string) (Result, error) {
	now := time.Now()

	c.mu.Lock()
	if e, ok := c.entries[imageRef]; ok && now.Before(e.expires) {
		c.mu.Unlock()
		return e.result, nil
	}
	c.mu.Unlock()

	res, err := c.inner.Classify(imageRef)
	if err != nil {
		return res, err
	}

	c.mu.Lock()
	c.entries[imageRef] = entry{result: res, expires: now.Add(c.ttl)}
	c.mu.Unlock()

	return res, nil
}
