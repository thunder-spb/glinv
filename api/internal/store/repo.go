package store

import (
	"net/url"
	"restapi/internal/model"
)

type AgentRepository interface {
	Insert(model.Data) error

	CheckServerExistsByHostname(string) (bool, error)
	GetIDByHostname(string) (int, error)
	DeletePackages(int) error

	Update(int, model.Data) error

	GetAll() ([]*model.Server, error)
	GetByHostname(string) (*model.Server, error)
}

type HostRepository interface {
	GetAll() ([]*model.Host, error)
	GetHostVars(int) ([]*model.HVar, error)
	GetHostByEnv(string) ([]*model.Host, error)
	GetHostByFilters(string, url.Values, []string) ([]*model.Host, error)
}

type GroupRepository interface {
	GetAll() ([]*model.Group, error)
}

type ServiceRepository interface {
	GetServiceByEnv(string) ([]*model.Service, error)
	GetServiceByLocation(string, string) ([]*model.Service, error)
}
