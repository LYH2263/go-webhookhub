package dispatch

import (
	"sort"

	"github.com/LYH2263/go-webhookhub/internal/endpoint"
)

// SortByID 稳定顺序，便于测试对照。
func SortByID(recs []endpoint.Record) {
	sort.SliceStable(recs, func(i, j int) bool { return recs[i].ID < recs[j].ID })
}

// StableJobs 按 EndpointID 排序。
func StableJobs(jobs []Job) {
	sort.SliceStable(jobs, func(i, j int) bool {
		if jobs[i].EndpointID == jobs[j].EndpointID {
			return jobs[i].DeliveryID < jobs[j].DeliveryID
		}
		return jobs[i].EndpointID < jobs[j].EndpointID
	})
}
