package webhookhub

// Close 标记关闭、先 Flush/Sync 日志再释放底层文件，最后把 HTTP 客户端置空。
// 之后 Dispatch / Subscribe 返回 ErrClosed，不会解引用 nil client。
func (h *Hub) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.closed {
		return nil
	}
	h.closed = true

	var first error
	if h.log != nil {
		if err := h.log.Flush(); err != nil && first == nil {
			first = err
		}
		if err := h.log.Sync(); err != nil && first == nil {
			first = err
		}
		if err := h.log.Close(); err != nil && first == nil {
			first = err
		}
	}
	h.client = nil
	h.transport = nil
	return first
}

// Closed 报告 Hub 是否已关闭。
func (h *Hub) Closed() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.closed
}
