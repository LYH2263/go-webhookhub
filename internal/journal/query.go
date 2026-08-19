package journal

func (l *Log) Seq() uint64 {
	if l == nil {
		return 0
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.seq
}

func (l *Log) SuccessCount() int {
	if l == nil {
		return 0
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	n := 0
	for _, r := range l.recs {
		if r.OK {
			n++
		}
	}
	return n
}

func (l *Log) FailCount() int {
	if l == nil {
		return 0
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	n := 0
	for _, r := range l.recs {
		if !r.OK {
			n++
		}
	}
	return n
}

func (l *Log) FilterEvent(event string, n int) []Record {
	all := l.Recent(0)
	out := make([]Record, 0)
	for i := len(all) - 1; i >= 0; i-- {
		if all[i].Event == event {
			out = append(out, all[i])
			if n > 0 && len(out) >= n {
				break
			}
		}
	}
	return out
}
