// Author: n0lsec
// Version: 0.1

package main

import (
	"flag"
	"fmt"
	"net"
)

func main() {
	ipRangePtr := flag.String("i", "", "IP range in CIDR notation (e.g., 104.36.248.0/22)")
	helpPtr := flag.Bool("h", false, "Show help")

	flag.Parse()

	if *helpPtr {
		printUsage()
		return
	}

	if *ipRangePtr == "" {
		printUsage()
		return
	}

	ips, err := getIPsInRange(*ipRangePtr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for ip := range ips {
		fmt.Println(ip)
	}
}

func getIPsInRange(cidrRange string) (<-chan string, error) {
	_, ipnet, err := net.ParseCIDR(cidrRange)
	if err != nil {
		return nil, err
	}

	ips := make(chan string, 1024)
	go func() {
		defer close(ips)
		for ip := ipnet.IP; ipnet.Contains(ip); incrementIP(ip) {
			ips <- ip.String()
		}
		// Remove the network and broadcast addresses
		<-ips
		<-ips
	}()

	return ips, nil
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  iprenger -i <IP range in CIDR notation>")
	fmt.Println("  iprenger -h (show help)")
}
