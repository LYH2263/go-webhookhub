package endpoint

import (
	"sort"
	"sync"

	ierr "github.com/LYH2263/go-webhookhub/internal/errors"
)

// Registry 内存端点表。所有取出的 Record 均为拷贝。
type Registry struct {
	mu       sync.Mutex
	byID     map[string]Record
	order    []string
	max      int
	allowHTTP bool
}

func NewRegistry(max int, allowHTTP bool) *Registry {
	if max <= 0 {
		max = 4096
	}
	return &Registry{
		byID:      make(map[string]Record),
		max:       max,
		allowHTTP: allowHTTP,
	}
}

func (r *Registry) Add(rec Record) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec = CloneRecord(rec)
	rec.ID = NormalizeID(rec.ID)
	rec.Events = NormalizeEvents(rec.Events)
	if err := Validate(rec, r.allowHTTP); err != nil {
		return "", err
	}
	if _, exists := r.byID[rec.ID]; exists {
		return "", ierr.Wrap(ierr.ErrAlreadyExists, rec.ID)
	}
	if len(r.byID) >= r.max {
		return "", ierr.ErrTooManyEndpoints
	}
	r.byID[rec.ID] = rec
	r.order = append(r.order, rec.ID)
	return rec.ID, nil
}

func (r *Registry) Restore(rec Record) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec = CloneRecord(rec)
	if rec.ID == "" {
		return ierr.Wrap(ierr.ErrInvalidEndpoint, "empty id")
	}
	if _, exists := r.byID[rec.ID]; exists {
		return ierr.Wrap(ierr.ErrAlreadyExists, rec.ID)
	}
	r.byID[rec.ID] = rec
	r.order = append(r.order, rec.ID)
	return nil
}

func (r *Registry) Remove(id string) error {
	_, err := r.Take(id)
	return err
}

func (r *Registry) Take(id string) (Record, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, ok := r.byID[id]
	if !ok {
		return Record{}, ierr.Wrap(ierr.ErrNotFound, id)
	}
	delete(r.byID, id)
	for i, x := range r.order {
		if x == id {
			r.order = append(r.order[:i], r.order[i+1:]...)
			break
		}
	}
	return CloneRecord(rec), nil
}

func (r *Registry) Get(id string) (Record, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, ok := r.byID[id]
	if !ok {
		return Record{}, ierr.Wrap(ierr.ErrNotFound, id)
	}
	return CloneRecord(rec), nil
}

func (r *Registry) SetEnabled(id string, enabled bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	rec, ok := r.byID[id]
	if !ok {
		return ierr.Wrap(ierr.ErrNotFound, id)
	}
	rec.Enabled = enabled
	r.byID[id] = rec
	return nil
}

func (r *Registry) List() []Record {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Record, 0, len(r.order))
	for _, id := range r.order {
		if rec, ok := r.byID[id]; ok {
			out = append(out, CloneRecord(rec))
		}
	}
	return out
}

func (r *Registry) Match(event string) []Record {
	return FilterByEvent(r.List(), event)
}

func (r *Registry) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.byID)
}

func (r *Registry) EnabledCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	n := 0
	for _, rec := range r.byID {
		if rec.Enabled {
			n++
		}
	}
	return n
}

func (r *Registry) IDs() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := append([]string(nil), r.order...)
	sort.Strings(append([]string(nil), out...))
	return out
}
