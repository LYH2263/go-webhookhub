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
		// 先落盘再释放：Flush 把 bufio 缓冲写到底层文件，Sync 落盘，
		// 最后 Close 关闭句柄。顺序反了（先 Close）会把 writer 置 nil、
		// 关闭文件，缓冲里最后几条记录就此丢失。
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
	if h.client != nil {
		h.client.CloseIdle()
	}
	h.client = nil
	return first
}

// Closed 报告 Hub 是否已关闭。
func (h *Hub) Closed() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.closed
}
