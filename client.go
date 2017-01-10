/*
Implements the solution to assignment 1 for UBC CS 416 2016 W2.

Usage:
$ go run client.go [local UDP ip:port] [server UDP ip:port]

Example:
$ go run client.go 192.168.0.16:2020 198.162.52.206:4116

*/

package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Marshall(guess uint32) ([]byte, error) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(guess)
	return network.Bytes(), err
}

func Average(low, high uint32) uint32 {
	return low + (high-low)/2
}

// Main workhorse method.
func main() {
	args := os.Args[1:]

	// // Missing command line args.
	if len(args) != 2 {
		fmt.Println("Usage: client.go [local UDP ip:port] [server UDP ip:port]")
		return
	}

	// Extract the command line args.
	local_ip_port := args[0]
	remote_ip_port := args[1]

	// Set up UDP connection
	LocalUDPAddr, err := net.ResolveUDPAddr("udp", local_ip_port)
	CheckError(err)

	RemoteUDPAddr, err := net.ResolveUDPAddr("udp", remote_ip_port)
	CheckError(err)

	connection, err := net.DialUDP("udp", LocalUDPAddr, RemoteUDPAddr)
	CheckError(err)

	defer connection.Close()

	ask := func(guess uint32) string {

		timeoutDuration, err := time.ParseDuration("3s")
		CheckError(err)

		message, err := Marshall(guess)
		CheckError(err)

		response := ""

		for {
			err = connection.SetWriteDeadline(time.Now().Add(timeoutDuration))
			CheckError(err)
			_, err = connection.Write(message)

			if err == nil {
				err = connection.SetReadDeadline(time.Now().Add(timeoutDuration))
				CheckError(err)

				buf := make([]byte, 1024)
				n, _, err := connection.ReadFromUDP(buf)
				if err == nil {
					response = string(buf[:n])
					break
				}
			}
		}
		return response
	}

	// Start searching
	var floor, ceil uint32 = 0, 4294967295

	// start the binary search

	for {
		guess := Average(floor, ceil)
		response := ask(guess)

		if strings.EqualFold(response, "low") {
			floor = guess
		} else if strings.EqualFold(response, "high") {
			ceil = guess
		} else {
			fmt.Println(response)
			break
		}
	}
}
