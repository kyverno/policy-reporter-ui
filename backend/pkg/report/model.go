package report

import "time"

// Resource API Model
type Resource struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Name       string `json:"name"`
	Namespace  string `json:"namespace,omitempty"`
	UID        string `json:"uid"`
}

// Result API Model
type Result struct {
	Message           string    `json:"message"`
	Policy            string    `json:"policy"`
	Rule              string    `json:"rule"`
	Priority          string    `json:"priority"`
	Status            string    `json:"status"`
	Severity          string    `json:"severity,omitempty"`
	Category          string    `json:"category,omitempty"`
	Scored            bool      `json:"scored"`
	Resource          Resource  `json:"resource"`
	CreationTimestamp time.Time `json:"creationTimestamp"`
}
