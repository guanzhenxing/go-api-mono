package health

import "time"

// Status represents the status of a service component
type Status struct {
	Status    string `json:"status"`
	Component string `json:"component"`
	Message   string `json:"message,omitempty"`
}

// Response represents the health check response
type Response struct {
	Status     string    `json:"status"`
	Version    string    `json:"version"`
	Components []Status  `json:"components"`
	Timestamp  time.Time `json:"timestamp"`
	Uptime     string    `json:"uptime"`
}
