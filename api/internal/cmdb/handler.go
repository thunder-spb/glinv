package cmdb

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"restapi/internal/model"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

//
// Index
//
func (s *server) index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, "REST API CMDB v.0.0.1")
	}
}

//
// Servers
//
func (s *server) dataServer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data model.Data

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		hostname := data["hard"]["hostname"]
		// Check the existence of an IP record and return the its ID
		exists, err := s.store.Agent().CheckServerExistsByHostname(hostname)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		if !exists {
			err = s.store.Agent().Insert(data)
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}
			log.Println("agent reports about new server:", hostname)

		} else {
			id, err := s.store.Agent().GetIDByHostname(hostname) // get the ID of an existing hostname
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			// Update server, packages and sysctl
			err = s.store.Agent().Update(id, data)
			if err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}

			log.Println("agent reports updated data about the server:", hostname)
		}

		s.respond(w, r, 200, nil)
	}
}

func (s *server) showServers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, "OK")
	}
}

func (s *server) showServerByHost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		server, err := s.store.Agent().GetByHostname(params["hostname"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, server)
	}
}

//
// Hosts
//
func (s *server) showHost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, "Host")
	}
}

func (s *server) showHosts() http.HandlerFunc {
	connection := "ssh"
	user := "pavel"
	port := 22

	return func(w http.ResponseWriter, r *http.Request) {
		nodes, err := s.store.Host().GetAll()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		// Map Construct
		mapInventoryHosts := make(map[string]interface{})
		mapEnv := make(map[string]interface{}) // Env

		mapGloblaVars := make(map[string]interface{}) // Global Vars
		mapGloblaVars["ansible_connection"] = connection
		mapGloblaVars["ansible_user"] = user
		mapGloblaVars["ansible_port"] = port
		mapEnv["vars"] = mapGloblaVars

		mapHosts := make(map[string]interface{}) // Hosts

		for _, node := range nodes {
			mapVars := make(map[string]interface{}) // Vars

			// Vars obligatory
			mapVars["ansible_host"] = node.IP
			//mapVars["ansible_environment"] = node.Environment

			// Vars optional
			nodeVars, err := s.store.Host().GetHostVars(node.ID)
			if err != nil {
				return
			}

			for _, v := range nodeVars {
				title := "ansible_" + v.Name

				if i, err := strconv.Atoi(v.Value); err == nil {
					mapVars[title] = i
				} else {
					mapVars[title] = v.Value
				}
			}

			mapHosts[node.Hostname] = mapVars
		}

		mapEnv["hosts"] = mapHosts
		mapInventoryHosts["all"] = mapEnv

		s.respond(w, r, http.StatusOK, mapInventoryHosts)
	}
}

func (s *server) showHostsByEnv() http.HandlerFunc {
	connection := "ssh"
	user := "pavel"
	port := 22

	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		// Get Nodes
		nodes, err := s.store.Host().GetHostByEnv(params["env"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		// Map Construct
		mapInventoryHosts := make(map[string]interface{})
		mapEnv := make(map[string]interface{}) // Env

		mapGloblaVars := make(map[string]interface{}) // Global Vars
		mapGloblaVars["ansible_connection"] = connection
		mapGloblaVars["ansible_user"] = user
		mapGloblaVars["ansible_port"] = port
		mapEnv["vars"] = mapGloblaVars

		mapHosts := make(map[string]interface{}) // Hosts

		for _, node := range nodes {
			mapVars := make(map[string]interface{}) // Vars

			// Vars obligatory
			mapVars["ansible_host"] = node.IP
			//mapVars["ansible_environment"] = node.Environment

			// Vars optional
			nodeVars, err := s.store.Host().GetHostVars(node.ID)
			if err != nil {
				return
			}

			for _, v := range nodeVars {
				title := "ansible_" + v.Name

				if i, err := strconv.Atoi(v.Value); err == nil {
					mapVars[title] = i
				} else {
					mapVars[title] = v.Value
				}
			}

			mapHosts[node.Hostname] = mapVars
		}

		mapEnv["hosts"] = mapHosts
		// var e string
		// if params["env"] == "prd" {
		// 	e = "production"
		// }
		// if params["env"] == "ppr" {
		// 	e = "preproduction"
		// }
		mapInventoryHosts[params["env"]] = mapEnv

		// Respond
		s.respond(w, r, http.StatusOK, mapInventoryHosts)
	}
}

func (s *server) showHostsByFilters() http.HandlerFunc {
	connection := "ssh"
	user := "pavel"
	port := 22

	return func(w http.ResponseWriter, r *http.Request) {
		mapParams := mux.Vars(r)
		env := "all"
		if len(mapParams["env"]) != 0 {
			env = mapParams["env"]
		}

		urlQuery := r.URL.Query()

		// Sorting parameters as is
		u, err := url.Parse(r.URL.String())
		if err != nil {
			panic(err)
		}

		var sliceSplitRawQuery []string
		parameters := strings.Split(u.RawQuery, "&")
		for _, parameter := range parameters {
			parts := strings.Split(parameter, "=")
			sliceSplitRawQuery = append(sliceSplitRawQuery, parts[0])
		}

		nodes := []*model.Host{}

		// Filter
		// Type Vars: tags, components, services, applications, roles
		filteredHosts := make(map[string]interface{})
		filters := []string{"tag", "component", "service", "application", "role"}

		// Check erros in request
		for key := range urlQuery {
			if !Contains(filters, key) {
				filteredHosts["ERROR"] = fmt.Sprintf("Invalid request parameter %v", key)
			}
		}

		nodes, _ = s.store.Host().GetHostByFilters(env, urlQuery, sliceSplitRawQuery)
		if len(nodes) != 0 {

			// Map Construct
			mapEnv := make(map[string]interface{}) // Env

			mapGloblaVars := make(map[string]interface{}) // Global Vars
			mapGloblaVars["ansible_connection"] = connection
			mapGloblaVars["ansible_user"] = user
			mapGloblaVars["ansible_port"] = port
			mapEnv["vars"] = mapGloblaVars

			mapHosts := make(map[string]interface{}) // Hosts

			for _, node := range nodes {
				mapVars := make(map[string]interface{}) // Vars

				// Vars obligatory
				mapVars["ansible_host"] = node.IP
				//mapVars["ansible_environment"] = node.Environment

				// Vars optional
				nodeVars, err := s.store.Host().GetHostVars(node.ID)
				if err != nil {
					return
				}

				for _, v := range nodeVars {
					title := "ansible_" + v.Name

					if i, err := strconv.Atoi(v.Value); err == nil {
						mapVars[title] = i
					} else {
						mapVars[title] = v.Value
					}
				}

				mapHosts[node.Hostname] = mapVars
			}

			mapEnv["hosts"] = mapHosts
			// var e string
			// if env == "prd" {
			// 	e = "production"
			// }
			// if env == "ppr" {
			// 	e = "preproduction"
			// }

			filteredHosts[env] = mapEnv

		} else {
			filteredHosts["FUCK"] = "There's a hell of a lot of errors in your request"
		}

		s.respond(w, r, http.StatusOK, filteredHosts)
	}
}

//
// Groups
//
func (s *server) showGroups() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		groups, err := s.store.Group().GetAll()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, groups)
	}
}

//
// Services
//
func (s *server) showServiceDesc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		description := `/services/prd
/services/ppr
/services/edu
/services/qa2
/services/qa
/services/dev2
/services/dev`

		s.respond(w, r, http.StatusOK, description)
	}
}

func (s *server) showServicesByEnv() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/x-yaml")

		params := mux.Vars(r)

		services, err := s.store.Service().GetServiceByEnv(params["env"])
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		// Map Construct
		mapInventory := make(map[string]interface{})
		mapEnv := make(map[string]interface{})
		mapLocation := make(map[string]interface{})

		for _, service := range services {
			mapDefault := make(map[string]interface{})
			mapService := make(map[string]interface{})

			items, err := s.store.Service().GetServiceByLocation(params["env"], service.Location)
			if err != nil {
				return
			}

			for _, item := range items {
				mapService[item.Title] = fmt.Sprintf("%s:%d", item.Value, item.Port)
			}

			mapDefault["default_node"] = mapService
			mapLocation[service.Location] = mapDefault
			mapEnv[service.Environment] = mapLocation
		}

		mapInventory["environment"] = mapEnv

		s.respond(w, r, http.StatusOK, mapInventory)
	}
}

//
// SSH
//

// showSSHConfig ...
func (s *server) showSSHConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sshConfig := make(map[string]interface{}) // Hosts

		allHosts, err := s.store.Host().GetAll()
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		for _, host := range allHosts {
			sshConfig[host.Hostname] = host.IP
		}

		s.sshConfigRespond(w, r, http.StatusOK, sshConfig)
	}

}

// showSSHConfig ...
func (s *server) showNginxConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nginxConfig := make(map[string]interface{})

		nginxConfig["test glinv"] = "glinv"

		s.nginxConfigRespond(w, r, http.StatusOK, nginxConfig)
	}

}
