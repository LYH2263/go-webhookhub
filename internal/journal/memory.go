package journal

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"

	"github.com/LYH2263/go-webhookhub/internal/clock"
	ierr "github.com/LYH2263/go-webhookhub/internal/errors"
)

// Log 环形内存日志，可选追加到文件。Close 前必须 Flush + Sync。
type Log struct {
	mu   sync.Mutex
	recs []Record
	cap  int
	seq  uint64
	clk  clock.Clock

	path   string
	file   *os.File
	writer *bufio.Writer
	closed bool
}

func New(capacity int, clk clock.Clock, path string) *Log {
	if capacity < 1 {
		capacity = 1
	}
	if clk == nil {
		clk = clock.Real{}
	}
	l := &Log{cap: capacity, clk: clk, recs: make([]Record, 0, min(capacity, 64))}
	if path != "" {
		l.path = path
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err == nil {
			l.file = f
			l.writer = bufio.NewWriterSize(f, 32*1024)
		}
	}
	return l
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (l *Log) Append(r Record) {
	if l == nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.closed {
		return
	}
	l.seq++
	r = CloneRecord(r)
	if r.At.IsZero() && l.clk != nil {
		r.At = l.clk.Now()
	}
	if len(l.recs) >= l.cap {
		n := copy(l.recs, l.recs[1:])
		l.recs = l.recs[:n]
		l.recs = append(l.recs, r)
	} else {
		l.recs = append(l.recs, r)
	}
	if l.writer != nil {
		enc, err := json.Marshal(r)
		if err == nil {
			_, _ = l.writer.Write(enc)
			_ = l.writer.WriteByte('\n')
		}
	}
}

func (l *Log) Recent(n int) []Record {
	if l == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if n <= 0 || n > len(l.recs) {
		n = len(l.recs)
	}
	start := len(l.recs) - n
	return CloneRecords(l.recs[start:])
}

func (l *Log) Len() int {
	if l == nil {
		return 0
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.recs)
}

// Flush 把缓冲写到文件，不关闭句柄。
func (l *Log) Flush() error {
	if l == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.writer == nil {
		return nil
	}
	if err := l.writer.Flush(); err != nil {
		return ierr.WrapErr(ierr.ErrFlush, err)
	}
	return nil
}

// Sync fsync 底层文件。调用前应 Flush。
func (l *Log) Sync() error {
	if l == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.writer != nil {
		if err := l.writer.Flush(); err != nil {
			return ierr.WrapErr(ierr.ErrFlush, err)
		}
	}
	if l.file == nil {
		return nil
	}
	if err := l.file.Sync(); err != nil {
		return ierr.WrapErr(ierr.ErrSync, err)
	}
	return nil
}

// Close 关闭底层文件句柄。调用前必须先 Flush+Sync，否则 bufio 缓冲里尚未
// 落盘的记录会随 writer 置 nil 而丢失。之后 Append 被忽略。
func (l *Log) Close() error {
	if l == nil {
		return nil
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.closed {
		return nil
	}
	l.closed = true
	l.writer = nil
	if l.file != nil {
		_ = l.file.Close()
		l.file = nil
	}
	return nil
}
