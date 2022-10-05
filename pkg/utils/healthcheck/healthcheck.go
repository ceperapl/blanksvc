package healthcheck

import (
	"net/http"
	"sync"
)

type Check func() error

// ...
type HealthChecker interface {
	AddLivenessChecks(check Check)
	AddReadinessChecks(check Check)
	LivenessHandler() http.HandlerFunc
	ReadinessHandler() http.HandlerFunc
}

func New() HealthChecker {
	return &healthcheck{}
}

type healthcheck struct {
	lock sync.RWMutex

	livenessChecks  []Check
	readinessChecks []Check
}

func (h *healthcheck) LivenessHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if err := h.checkLiveness(); err != nil {
			rw.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		rw.WriteHeader(http.StatusOK)
	}
}

func (h *healthcheck) ReadinessHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if err := h.checkReadiness(); err != nil {
			rw.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		rw.WriteHeader(http.StatusOK)
	}
}

func (h *healthcheck) AddLivenessChecks(check Check) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.livenessChecks = append(h.livenessChecks, check)
}

func (h *healthcheck) AddReadinessChecks(check Check) {
	h.lock.Lock()
	defer h.lock.Unlock()

	h.readinessChecks = append(h.readinessChecks, check)
}

func (h *healthcheck) checkReadiness() error {
	h.lock.RLock()
	defer h.lock.RUnlock()

	for _, check := range h.readinessChecks {
		if err := check(); err != nil {
			return err
		}
	}
	return nil
}

func (h *healthcheck) checkLiveness() error {
	h.lock.RLock()
	defer h.lock.RUnlock()

	for _, check := range h.livenessChecks {
		if err := check(); err != nil {
			return err
		}
	}
	return nil
}
