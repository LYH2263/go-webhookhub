package dispatch

import "github.com/LYH2263/go-webhookhub/internal/endpoint"

// DropDisabled 从计划源去掉未启用端点。
func DropDisabled(in []endpoint.Record) []endpoint.Record {
	out := make([]endpoint.Record, 0, len(in))
	for _, r := range in {
		if r.Enabled {
			out = append(out, endpoint.CloneRecord(r))
		}
	}
	return out
}

func EmptyJobs() []Job { return []Job{} }
