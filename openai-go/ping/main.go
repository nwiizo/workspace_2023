package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

const (
	red   = "\033[31m"
	blue  = "\033[34m"
	reset = "\033[0m"
)

func main() {
	hosts, err := readHostsFromFile("hosts.txt")
	if err != nil {
		fmt.Printf("Failed to read hosts from file: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(hosts))

	for _, host := range hosts {
		go func(host string) {
			for {
				ping(host)
				time.Sleep(1 * time.Second)
			}
			wg.Done()
		}(host)
	}

	wg.Wait()
}

func readHostsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var hosts []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hosts = append(hosts, scanner.Text())
	}

	return hosts, scanner.Err()
}

func ping(host string) {
	timeout := time.Duration(1 * time.Second)
	conn, err := net.DialTimeout("ip4:icmp", host, timeout)

	if err != nil {
		fmt.Printf("%sRequest timeout for ICMP packet to %s: %v%s\n", red, host, err, reset)
		return
	}
	defer conn.Close()

	start := time.Now()
	conn.SetDeadline(start.Add(timeout))

	err = sendICMPEchoRequest(conn)
	if err != nil {
		fmt.Printf("%sFailed to send ICMP echo request to %s: %v%s\n", red, host, err, reset)
		return
	}

	err = receiveICMPEchoReply(conn)
	if err != nil {
		fmt.Printf("%sFailed to receive ICMP echo reply from %s: %v%s\n", red, host, err, reset)
		return
	}

	elapsed := time.Since(start)

	fmt.Printf("%sPing to %s: success, time=%v%s\n", blue, host, elapsed, reset)
}

func sendICMPEchoRequest(conn net.Conn) error {
	// Create ICMP echo request packet
	packet := make([]byte, 8)
	packet[0] = 8 // ICMP type: echo request
	packet[1] = 0 // ICMP code: no code for echo request
	packet[2] = 0 // Checksum (initially 0)
	packet[3] = 0 // Checksum (initially 0)
	packet[4] = 0 // Identifier
	packet[5] = 0 // Identifier
	packet[6] = 0 // Sequence number
	packet[7] = 1 // Sequence number

	// Calculate the checksum
	checksum := calculateChecksum(packet)
	packet[2] = byte(checksum >> 8)
	packet[3] = byte(checksum & 0xff)

	// Send the packet
	_, err := conn.Write(packet)
	return err
}

func receiveICMPEchoReply(conn net.Conn) error {
	reply := make([]byte, 1024)

	// Receive the packet
	_, err := conn.Read(reply)
	return err
}

func calculateChecksum(data []byte) uint16 {
	var sum uint32

	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(data[i])<<8 | uint32(data[i+1])
	}

	if len(data)%2 == 1 {
		sum += uint32(data[len(data)-1]) << 8
	}

	for sum>>16 > 0 {
		sum = (sum & 0xffff) + (sum >> 16)
	}

	return ^uint16(sum)
}
