package journal

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/LYH2263/go-webhookhub/internal/endpoint"
	ierr "github.com/LYH2263/go-webhookhub/internal/errors"
	"github.com/LYH2263/go-webhookhub/internal/payload"
)

type endpointFile struct {
	Endpoints []endpoint.Snapshot `json:"endpoints"`
}

// FileSink 把端点快照原子写入 JSON 文件。
type FileSink struct {
	Path string
}

func NewFileSink(path string) *FileSink {
	return &FileSink{Path: path}
}

func (s *FileSink) SaveEndpoints(views []endpointView) error {
	if s == nil || s.Path == "" {
		return nil
	}
	snaps := make([]endpoint.Snapshot, 0, len(views))
	for _, v := range views {
		snaps = append(snaps, endpoint.Snapshot{
			ID:          v.ID,
			URL:         v.URL,
			Secret:      payload.CloneBytes(v.Secret),
			Events:      payload.CloneStrings(v.Events),
			Enabled:     v.Enabled,
			Headers:     payload.CloneStringMap(v.Headers),
			TimeoutNS:   int64(v.TimeoutNS),
			MaxAttempts: v.MaxAttempts,
			Description: v.Description,
		})
	}
	raw, err := json.MarshalIndent(endpointFile{Endpoints: snaps}, "", "  ")
	if err != nil {
		return ierr.WrapErr(ierr.ErrPersist, err)
	}
	return AtomicWriteFile(s.Path, raw)
}

// endpointView 与库面 EndpointView 字段对齐的最小结构，避免 journal 导入根包。
type endpointView struct {
	ID          string
	URL         string
	Events      []string
	Enabled     bool
	Secret      []byte
	Headers     map[string]string
	TimeoutNS   int64
	MaxAttempts int
	Description string
}

// SaveJSON 供根包通过适配器调用。Hub 传入 JSON 可序列化视图。
func (s *FileSink) SaveJSON(views []EndpointJSON) error {
	if s == nil || s.Path == "" {
		return nil
	}
	snaps := make([]endpoint.Snapshot, 0, len(views))
	for _, v := range views {
		snaps = append(snaps, endpoint.Snapshot{
			ID:          v.ID,
			URL:         v.URL,
			Secret:      payload.CloneBytes(v.Secret),
			Events:      payload.CloneStrings(v.Events),
			Enabled:     v.Enabled,
			Headers:     payload.CloneStringMap(v.Headers),
			TimeoutNS:   v.TimeoutNS,
			MaxAttempts: v.MaxAttempts,
			Description: v.Description,
		})
	}
	raw, err := json.MarshalIndent(endpointFile{Endpoints: snaps}, "", "  ")
	if err != nil {
		return ierr.WrapErr(ierr.ErrPersist, err)
	}
	if err := AtomicWriteFile(s.Path, raw); err != nil {
		return ierr.WrapErr(ierr.ErrPersist, err)
	}
	return nil
}

type EndpointJSON struct {
	ID          string            `json:"id"`
	URL         string            `json:"url"`
	Secret      []byte            `json:"secret"`
	Events      []string          `json:"events"`
	Enabled     bool              `json:"enabled"`
	Headers     map[string]string `json:"headers,omitempty"`
	TimeoutNS   int64             `json:"timeout_ns"`
	MaxAttempts int               `json:"max_attempts"`
	Description string            `json:"description,omitempty"`
}

func LoadEndpointFile(path string) ([]endpoint.Snapshot, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		return nil, ierr.WrapErr(ierr.ErrPersist, err)
	}
	var blob endpointFile
	if err := json.Unmarshal(raw, &blob); err != nil {
		return nil, ierr.WrapErr(ierr.ErrPersist, err)
	}
	out := make([]endpoint.Snapshot, 0, len(blob.Endpoints))
	for _, s := range blob.Endpoints {
		s.Secret = payload.CloneBytes(s.Secret)
		s.Events = payload.CloneStrings(s.Events)
		s.Headers = payload.CloneStringMap(s.Headers)
		out = append(out, s)
	}
	return out, nil
}

func AtomicWriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if dir == "" {
		dir = "."
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".webhookhub-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	ok := false
	defer func() {
		if !ok {
			_ = os.Remove(tmpName)
		}
	}()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpName, path); err != nil {
		return err
	}
	ok = true
	return nil
}
