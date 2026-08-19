package webhookhub

// Stats 返回当前计数快照。
func (h *Hub) Stats() Stats {
	h.mu.Lock()
	defer h.mu.Unlock()
	enabled := 0
	n := 0
	if h.reg != nil {
		n = h.reg.Len()
		enabled = h.reg.EnabledCount()
	}
	jlen := 0
	if h.log != nil {
		jlen = h.log.Len()
	}
	return Stats{
		Endpoints:     n,
		Enabled:       enabled,
		Dispatches:    h.dispatches,
		Deliveries:    h.deliveries,
		Successes:     h.successes,
		Failures:      h.failures,
		LastEvent:     h.lastEvent,
		LastAt:        h.lastAt,
		JournalLen:    jlen,
		Closed:        h.closed,
		DefaultMaxTry: h.maxAttempts,
	}
}

// EndpointCount 已注册端点数。
func (h *Hub) EndpointCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.reg == nil {
		return 0
	}
	return h.reg.Len()
}
