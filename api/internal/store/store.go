package store

// Store interface
type Store interface {
	Agent() AgentRepository
	Host() HostRepository
	Group() GroupRepository
	Service() ServiceRepository
}
