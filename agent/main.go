package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"time"

	sigar "github.com/cloudfoundry/gosigar"
	sysctl "github.com/lorenzosaino/go-sysctl"
	"github.com/zcalusic/sysinfo"
)

const glinvURL = "http://192.168.10.229:10011"

const refreshRate = 60 // sec

func checkURL(url string) int {
	log.Println("Check url GLINV", glinvURL, "Refresh rate:", refreshRate, "sec")

	resp, err := http.Get(url)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0
	}

	return 1
}

func checkLoop() {
	for {
		check := checkURL(glinvURL)

		if check == 1 {
			resource := "/agent"

			// sys
			var sys sysinfo.SysInfo
			sys.GetSysInfo()

			// ip
			ip, err := getIP()
			if err != nil {
				log.Println(err)
			}

			// hostname
			hostname, err := os.Hostname()
			if err != nil {
				log.Panicln(err)
			}

			// memory
			mem := sigar.Mem{}
			mem.Get()
			ram := fmt.Sprintf("%v\n%v\n%v", formatMemory(mem.Total), formatMemory(mem.Used), formatMemory(mem.Free))

			// hdd
			disk := DiskUsage("/")
			hdd := fmt.Sprintf("%v\n%v\n%v", float64(disk.All)/float64(GB), float64(disk.Used)/float64(GB), float64(disk.Free)/float64(GB))

			// uptime
			uptime := sigar.Uptime{}
			uptime.Get()

			// sysctl
			sysctlconf, err := sysctl.GetAll()

			data := map[string]map[string]string{
				"hard": {
					"hostname":      hostname,
					"ip":            ip,
					"osName":        sys.OS.Name,
					"osArch":        sys.OS.Architecture,
					"kernelRelease": sys.Kernel.Release,
					"modelCPU":      sys.CPU.Model,
					"numCPU":        strconv.Itoa(runtime.NumCPU()),
					"ram":           ram,
					"hdd":           hdd,
					"uptime":        uptime.Format(),
					"resolv.conf":   getResolv(),
					"mtu":           mtu(),
				},
				"pkg":    packages(sys.OS.Vendor),
				"sysctl": sysctlconf,
			}

			byteData, err := json.Marshal(data)
			if err != nil {
				log.Fatalln(err)
			}

			u, _ := url.ParseRequestURI(glinvURL)
			u.Path = resource
			urlStr := u.String()

			client := &http.Client{}
			r, _ := http.NewRequest("PUT", urlStr, bytes.NewBuffer(byteData))
			r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
			r.Header.Add("Content-Type", "application/json")
			resp, err := client.Do(r)
			if err != nil {
				log.Fatalln(err)
			}

			log.Println(hostname, ip, resp.Status)

		} else {
			log.Println("GLINV Server is not available, waiting availability...")
		}

		time.Sleep(refreshRate * time.Second)
	}
}

func main() {
	checkLoop()

}
