package model

// Service model
type Service struct {
	ID          int
	Environment string
	Location    string
	Hostname    string
	IP          string
	Title       string
	Value       string
	Port        int
	Approved    bool
}
