package model

// Server model
type Server struct {
	ID             int    `json:"id"`
	Hostname       string `json:"hostname"`
	IP             string `json:"ip"`
	OSName         string `json:"os_name"`
	OSVendor       string `json:"os_vendor"`
	OSVersion      string `json:"os_version"`
	OSRelease      string `json:"os_release"`
	OSArchitecture string `json:"os_architecture"`
	KernelRelease  string `json:"kernel_release"`
	CPUModel       string `json:"cpu_model "`
}

// Data ...
type Data map[string]map[string]string

// SysCtlServer model
type SysCtlServer struct {
	Name  string `db:"name"`
	Value string `db:"value"`
}
