package api

import (
	"net/http"

	"github.com/pulseguard/pulseguard/internal/server/store"
)

type dashboardResponse struct {
	TotalAgents    int           `json:"total_agents"`
	TotalJobs      int           `json:"total_jobs"`
	SuccessRate    float64       `json:"success_rate"`
	RecentFailures []interface{} `json:"recent_failures"`
}

func dashboardHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		agentCount, _ := s.AgentCount()
		jobCount, _ := s.JobCount()
		successRate, _ := s.SuccessRate()

		failures, _ := s.RecentFailures(10)
		failureList := make([]interface{}, 0, len(failures))
		for _, f := range failures {
			failureList = append(failureList, f)
		}

		writeJSON(w, http.StatusOK, dashboardResponse{
			TotalAgents:    agentCount,
			TotalJobs:      jobCount,
			SuccessRate:    successRate,
			RecentFailures: failureList,
		})
	}
}
