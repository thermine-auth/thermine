package activity

import (
	"time"

	"xermess/internal/model"
	"xermess/internal/store"
)

// eventResponse is one line of the activity list on the dashboard.
type eventResponse struct {
	ID        string    `json:"id"`
	Action    string    `json:"action"`
	Actor     string    `json:"actor"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}

// logResponse is the same line on the logs page, which has room for more.
type logResponse struct {
	eventResponse
	UserAgent  string `json:"user_agent"`
	TargetType string `json:"target_type"`
}

// overviewResponse is the dashboard: what there is, and what just happened.
type overviewResponse struct {
	Counts   store.Counts    `json:"counts"`
	Activity []eventResponse `json:"activity"`
}

func newOverviewResponse(counts store.Counts, events []model.AuditLog) overviewResponse {
	return overviewResponse{Counts: counts, Activity: newEventResponses(events)}
}

func newEventResponses(events []model.AuditLog) []eventResponse {
	out := make([]eventResponse, 0, len(events))
	for _, event := range events {
		out = append(out, newEventResponse(event))
	}

	return out
}

func newEventResponse(event model.AuditLog) eventResponse {
	return eventResponse{
		ID:        event.ID.String(),
		Action:    event.Action,
		Actor:     event.ActorEmail,
		IP:        event.IP,
		CreatedAt: event.CreatedAt,
	}
}

func newLogResponses(events []model.AuditLog) []logResponse {
	out := make([]logResponse, 0, len(events))
	for _, event := range events {
		out = append(out, logResponse{
			eventResponse: newEventResponse(event),
			UserAgent:     event.UserAgent,
			TargetType:    event.TargetType,
		})
	}

	return out
}
