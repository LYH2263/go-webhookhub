package webhookhub

import (
	"github.com/LYH2263/go-webhookhub/internal/journal"
	"github.com/LYH2263/go-webhookhub/internal/payload"
)

type fileSinkAdapter struct {
	inner *journal.FileSink
}

func (a fileSinkAdapter) SaveEndpoints(views []EndpointView) error {
	if a.inner == nil {
		return nil
	}
	js := make([]journal.EndpointJSON, 0, len(views))
	for _, v := range views {
		js = append(js, journal.EndpointJSON{
			ID:          v.ID,
			URL:         v.URL,
			Secret:      payload.CloneBytes(v.Secret),
			Events:      payload.CloneStrings(v.Events),
			Enabled:     v.Enabled,
			Headers:     payload.CloneStringMap(v.Headers),
			TimeoutNS:   int64(v.Timeout),
			MaxAttempts: v.MaxAttempts,
			Description: v.Description,
		})
	}
	_ = a.inner.SaveJSON(js)
	return nil
}
