package webhookhub

import (
	"github.com/LYH2263/go-webhookhub/internal/endpoint"
	"github.com/LYH2263/go-webhookhub/internal/payload"
)

// Subscribe 注册端点。Secret/Events/Headers 入库前拷贝。
// 若配置了持久化且写入失败，会回滚本次注册并返回错误。
func (h *Hub) Subscribe(ep Endpoint) (string, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := h.checkOpenLocked(); err != nil {
		return "", err
	}
	stored := endpoint.Record{
		ID:          ep.ID,
		URL:         ep.URL,
		Secret:      payload.CloneBytes(ep.Secret),
		Events:      payload.CloneStrings(ep.Events),
		Enabled:     ep.Enabled,
		Headers:     payload.CloneStringMap(ep.Headers),
		Timeout:     ep.Timeout,
		MaxAttempts: ep.MaxAttempts,
		Description: ep.Description,
	}
	id, err := h.reg.Add(stored)
	if err != nil {
		return "", err
	}
	if h.sink != nil {
		if err := h.sink.SaveEndpoints(h.viewsLocked()); err != nil {
			_ = h.reg.Remove(id)
			return "", err
		}
	}
	return id, nil
}

// Unsubscribe 删除端点。持久化失败时把端点加回去（回滚删除）。
func (h *Hub) Unsubscribe(id string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := h.checkOpenLocked(); err != nil {
		return err
	}
	removed, err := h.reg.Take(id)
	if err != nil {
		return err
	}
	if h.sink != nil {
		if err := h.sink.SaveEndpoints(h.viewsLocked()); err != nil {
			_ = h.reg.Restore(removed)
			return err
		}
	}
	return nil
}

// ListEndpoints 返回端点视图；Secret 为独立拷贝。
func (h *Hub) ListEndpoints() []EndpointView {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.viewsLocked()
}

func (h *Hub) viewsLocked() []EndpointView {
	recs := h.reg.List()
	out := make([]EndpointView, 0, len(recs))
	for _, r := range recs {
		out = append(out, recordToView(r))
	}
	return out
}

func recordToView(r endpoint.Record) EndpointView {
	return EndpointView{
		ID:          r.ID,
		URL:         r.URL,
		Events:      payload.CloneStrings(r.Events),
		Enabled:     r.Enabled,
		SecretLen:   len(r.Secret),
		Secret:      r.Secret,
		Headers:     payload.CloneStringMap(r.Headers),
		Timeout:     r.Timeout,
		MaxAttempts: r.MaxAttempts,
		Description: r.Description,
	}
}

// SetEnabled 启用或禁用端点。
func (h *Hub) SetEnabled(id string, enabled bool) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if err := h.checkOpenLocked(); err != nil {
		return err
	}
	if err := h.reg.SetEnabled(id, enabled); err != nil {
		return err
	}
	if h.sink != nil {
		if err := h.sink.SaveEndpoints(h.viewsLocked()); err != nil {
			_ = h.reg.SetEnabled(id, !enabled)
			return err
		}
	}
	return nil
}

func (h *Hub) persistLocked() error {
	if h.sink == nil {
		return nil
	}
	return h.sink.SaveEndpoints(h.viewsLocked())
}
