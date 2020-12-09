package model

// Host model
type Host struct {
	ID          int    `yaml:"id" json:"id"`
	Hostname    string `yaml:"hostname" json:"hostname"`
	IP          string `yaml:"ip" json:"ip"`
	Environment string `yaml:"environment" json:"environment"`
}

// HVar model
type HVar struct {
	ID    int
	Name  string
	Value string
}
