package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

//
// IP
//

func getIP() (string, error) {
	ifaces, err := net.Interfaces()

	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		addrs, err := iface.Addrs()

		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()

			if ip == nil {
				continue // not an ipv4 address
			}

			return ip.String(), nil
		}
	}

	return "", errors.New("are you connected to the network?")
}

//
// resolv.conf
//

func getResolv() string {
	file, err := os.Open("/etc/resolv.conf")
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string
	re := regexp.MustCompile(`^nameserver.*|^options.*|^search.*`)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && re.FindStringSubmatch(line) != nil {
			txtlines = append(txtlines, scanner.Text())
		}
	}
	file.Close()
	resolvfile := strings.Join(txtlines, "\n")
	return resolvfile
}

//
// ram
//

// formatMemory ...
func formatMemory(val uint64) uint64 {
	return val / 1024
}

//
// hdd
//

// DiskStatus ...
type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

// DiskUsage ...
func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

//
// mtu
//

func mtu() string {
	var mtu string
	out, _ := exec.Command("cat", "/sys/class/net/ens192/mtu").Output()
	mtu = fmt.Sprintf("%s", out)

	if len(mtu) == 0 {
		mtu = "0"
	}

	return mtu
}
