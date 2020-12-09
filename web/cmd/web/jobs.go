package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"glinv/pkg/models"
	"log"
	"net/http"
	"strings"
	"time"
)

//
// Consul
//

// 60s
func secTicker() *time.Ticker {
	return time.NewTicker(time.Second * time.Duration(61-time.Now().Second()))
}

// 10m
func minuteTicker() *time.Ticker {
	return time.NewTicker(time.Second * time.Duration(600-time.Now().Second()))
}

// 12h
func hourTicker() *time.Ticker {
	return time.NewTicker(time.Hour * time.Duration(35-time.Now().Hour()))
}

// hJobs ...
func hJobs(db *sql.DB) error {
	h := hourTicker()
	for {
		<-h.C
		h = hourTicker()

		// JOB: VACUUM
		if _, err := db.Exec("VACUUM (FULL, VERBOSE) public.server_hard_agent;"); err != nil {
			return err
		}

		if _, err := db.Exec("VACUUM (FULL, VERBOSE) public.server_sysctl_agent;"); err != nil {
			return err
		}

	}
}

// mJobs ...
func mJobs(db *sql.DB) error {
	m := minuteTicker()
	for {
		<-m.C
		m = minuteTicker()

		// JOB: Compare servers with hosts
		// Select all servers
		rows, err := db.Query(`SELECT id, hostname, ip FROM server_hard_agent`)
		if err != nil {
			return err
		}
		defer rows.Close()

		srv := []*models.ServerAgent{}
		for rows.Next() {
			s := &models.ServerAgent{}
			err := rows.Scan(
				&s.ID,
				&s.Hostname,
				&s.IP,
			)

			if err != nil {
				return err
			}

			srv = append(srv, s)
		}

		if err = rows.Err(); err != nil {
			return err
		}

		// Select all hosts
		rows, err = db.Query(`SELECT id, hostname, ip FROM inventory_hosts`)
		if err != nil {
			return err
		}
		defer rows.Close()

		hst := []*models.InventoryHost{}

		for rows.Next() {
			h := &models.InventoryHost{}
			err := rows.Scan(
				&h.ID,
				&h.Hostname,
				&h.IP,
			)

			if err != nil {
				return err
			}

			hst = append(hst, h)
		}

		if err = rows.Err(); err != nil {
			return err
		}

		// Clear status in inventory_hosts and server_hard_agent
		stmtComapreDelete := `DELETE FROM compare_server_host;`
		if _, err := db.Exec(stmtComapreDelete); err != nil {
			return err
		}

		// Compare
		for _, s := range srv {
			for _, h := range hst {
				if (s.Hostname == h.Hostname) && (s.IP == h.IP) {
					if _, err := db.Exec("INSERT INTO compare_server_host VALUES($1, $2)", s.ID, h.ID); err != nil {
						return err
					}
				}
			}
		}

	}
}

// sJobs ...
func sJobs(db *sql.DB) error {
	t := secTicker()
	for {
		<-t.C
		t = secTicker()

		// JOB: Checking the server uptime and alerting
		rows, err := db.Query(`SELECT hostname, ip, uptime FROM server_hard_agent`)
		if err != nil {
			return err
		}
		defer rows.Close()

		servers := []*models.ServerAgent{}

		for rows.Next() {
			server := &models.ServerAgent{}
			err := rows.Scan(
				&server.Hostname,
				&server.IP,
				&server.Uptime,
			)

			if err != nil {
				return err
			}

			servers = append(servers, server)
		}

		if err = rows.Err(); err != nil {
			return err
		}

		for _, s := range servers {
			stmt := `SELECT alert FROM server_alerts WHERE server_id = $1 limit 1;`
			alert, _ := db.Query(stmt, s.ID)
			if !alert.Next() {
				uptime := strings.Split(s.Uptime, ":")
				if uptime[0] == " 0" && uptime[1] == "01" {
					url := "#"
					var bearer = "Bearer " + "8fd4cnc7ibbhzggnhykqifdhae"

					msg := fmt.Sprintf("##### [Alert](#) \n *** \n :warning: uptime  %v %v  less than 1 minute \n***", s.Hostname, s.IP)
					data := map[string]string{
						"channel_id": "phj1pztnz3833c9qbr841rdz1o",
						"message":    msg,
					}

					byteData, err := json.Marshal(data)
					if err != nil {
						log.Fatalln(err)
					}

					req, err := http.NewRequest("POST", url, bytes.NewBuffer(byteData))

					req.Header.Add("Authorization", bearer)
					req.Header.Set("X-Custom-Header", "glinv")
					req.Header.Set("Content-Type", "application/json")

					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						panic(err)
					}
					defer resp.Body.Close()
				}
			}
		}

		// // JOB: Get Services From Consul
		// urlS := fmt.Sprintf("#v1/catalog/services")
		// reqS, err := http.NewRequest("GET", urlS, nil)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// respS, err := http.DefaultClient.Do(reqS)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer respS.Body.Close()

		// bS, _ := ioutil.ReadAll(respS.Body)

		// var dataS map[string]interface{}
		// if err := json.Unmarshal(bS, &dataS); err != nil {
		// 	panic(err)
		// }

		// for service := range dataS {
		// 	tags := dataS[service].([]interface{})
		// 	for value := range tags {
		// 		tag := tags[value].(string)

		// 		if tag == "dev" {
		// 			//fmt.Println(service, tag)
		// 			if _, err := db.Exec("UPDATE inventory_services SET status_in_consul = $1 WHERE title = $2", true, service); err != nil {
		// 				return err
		// 			}
		// 		}
		// 	}

		// }

		// // JOB: Get Nodes From Consul
		// urlN := fmt.Sprintf("#/v1/catalog/nodes")
		// reqN, err := http.NewRequest("GET", urlN, nil)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// respN, err := http.DefaultClient.Do(reqN)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer respN.Body.Close()

		// bN, _ := ioutil.ReadAll(respN.Body)

		// var dataN []map[string]interface{}

		// if err := json.Unmarshal(bN, &dataN); err != nil {
		// 	panic(err)
		// }

		//log.Println(dataN)

	}
}
