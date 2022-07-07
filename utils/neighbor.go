package utils

import (
	"fmt"
	"gobc/def"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

func IsFoundHost(host string, port uint16) bool {
	target := fmt.Sprintf("%s:%d", host, port)

	_, err := net.DialTimeout("tcp", target, 1*time.Second)
	if err != nil {
		fmt.Printf("%s %v\n", target, err)
		return false
	}
	color.Green(target)
	return true
}

var PATTERN = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

func FindNeighbors(myHost string, myPort uint16, startIP uint8, endIP uint8, startPort uint16, endPort uint16) []string {
	address := fmt.Sprintf("%s:%d", myHost, myPort)

	m := PATTERN.FindStringSubmatch(myHost)
	if m == nil {
		return nil
	}

	prefixHost := m[1]
	lastIP, _ := strconv.Atoi(m[len(m)-1])
	neighbors := make([]string, 0)

	color.Cyan("\n" + strings.Repeat("=", 10) + "Serching for other nodes" + strings.Repeat("=", 10))
	for port := startPort; port <= endPort; port += 1 {
		for ip := startIP; ip <= endIP; ip += 1 {
			guessHost := fmt.Sprintf("%s%d", prefixHost, lastIP+int(ip))
			guessTarget := fmt.Sprintf("%s:%d", guessHost, port)
			if guessTarget != address && IsFoundHost(guessHost, port) {
				neighbors = append(neighbors, guessTarget)
			}
		}
	}
	color.Cyan(strings.Repeat("=", 10) + "Serching for other nodes finished" + strings.Repeat("=", 10) + "\n")
	return neighbors
}

//自身のIPアドレスを取得
func GetHost() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return def.SELF_IP
	}

	for _, addr := range addrs {
		ip, ok := addr.(*net.IPNet)
		if ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String()
			}
		}
	}
	return def.SELF_IP
}
