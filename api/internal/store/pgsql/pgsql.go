package pgsql

import (
	"database/sql"
	"restapi/internal/store"
)

// Store ...
type Store struct {
	db                *sql.DB
	agentRepository   *AgentRepository
	hostRepository    *HostRepository
	groupRepository   *GroupRepository
	serviceRepository *ServiceRepository
}

// New ...
func New(db *sql.DB) *Store {
	return &Store{db: db}
}

// Agent ...
func (s *Store) Agent() store.AgentRepository {
	if s.agentRepository != nil {
		return s.agentRepository
	}

	s.agentRepository = &AgentRepository{store: s}

	return s.agentRepository
}

// Host Store
func (s *Store) Host() store.HostRepository {
	if s.hostRepository != nil {
		return s.hostRepository
	}

	s.hostRepository = &HostRepository{store: s}

	return s.hostRepository
}

// Group Store
func (s *Store) Group() store.GroupRepository {
	if s.groupRepository != nil {
		return s.groupRepository
	}

	s.groupRepository = &GroupRepository{store: s}

	return s.groupRepository
}

// Service Store
func (s *Store) Service() store.ServiceRepository {
	if s.groupRepository != nil {
		return s.serviceRepository
	}

	s.serviceRepository = &ServiceRepository{store: s}

	return s.serviceRepository
}
