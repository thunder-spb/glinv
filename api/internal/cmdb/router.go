package cmdb

func (s *server) configureRouter() {
	s.router.HandleFunc("/", s.index()).Methods("GET")

	s.router.HandleFunc("/agent", s.dataServer()).Methods("PUT")
	s.router.HandleFunc("/servers", s.showServers()).Methods("GET")
	s.router.HandleFunc("/servers/{hostname}", s.showServerByHost()).Methods("GET")

	s.router.HandleFunc("/hosts", s.showHosts()).Methods("GET")
	s.router.HandleFunc("/hosts/filter", s.showHostsByFilters()).Methods("GET")
	s.router.HandleFunc("/hosts/{env}", s.showHostsByEnv()).Methods("GET")
	s.router.HandleFunc("/hosts/{env}/filter", s.showHostsByFilters()).Methods("GET")
	s.router.HandleFunc("/hosts/{hostname}", s.showHost()).Methods("GET") //TODO

	s.router.HandleFunc("/groups", s.showGroups()).Methods("GET")

	s.router.HandleFunc("/service", s.showServiceDesc()).Methods("GET")
	s.router.HandleFunc("/services/{env}", s.showServicesByEnv()).Methods("GET")

	s.router.HandleFunc("/ssh/config", s.showSSHConfig()).Methods("GET")
	s.router.HandleFunc("/nginx/config/test", s.showNginxConfig()).Methods("GET")
}
