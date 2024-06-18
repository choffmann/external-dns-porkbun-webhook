package entities

import "sync"

type HealthStatus struct {
	health bool
	ready  bool
	rwLock sync.RWMutex
}

func NewHealthStatus() *HealthStatus {
	return &HealthStatus{
		health: false,
		ready:  false,
		rwLock: sync.RWMutex{},
	}
}

func (h *HealthStatus) GetReady() bool {
	h.rwLock.RLock()
	ready := h.ready
	h.rwLock.RUnlock()

	return ready
}

func (h *HealthStatus) GetHealth() bool {
	h.rwLock.RLock()
	health := h.health
	h.rwLock.RUnlock()

	return health
}

func (h *HealthStatus) SetReady(ready bool) {
	h.rwLock.Lock()
	h.ready = ready
	h.rwLock.Unlock()
}

func (h *HealthStatus) SetHealth(healthy bool) {
	h.rwLock.Lock()
	h.health = healthy
	h.rwLock.Unlock()
}
