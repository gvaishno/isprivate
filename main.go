/*
Author : Vaishno Chaitanya
License : GNU GPL v3
*/

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"unicode"
)

// Is IP address valid or not
func isValidIp(ip *string) bool {

	var validIP bool

	ipaddr := net.ParseIP(*ip)

	if ipaddr == nil {
		validIP = false
	} else {
		validIP = true
	}

	return validIP
}

// Check if the IP address is in the range of the subnet
func isInRange(ip *string, subnet *string) bool {
	var inRange bool

	x := *ip

	y := *subnet

	_, ipnetA, err := net.ParseCIDR(y)

	if err != nil {
		log.Fatal(err)
	}

	ipB := net.ParseIP(x)

	if ipnetA.Contains(ipB) {
		inRange = true
	}

	return inRange
}

// Parse private_blocks.csv
func parsePrivateBlocks() []string {
	var privateBlocks []string

	f, err := os.Open("private_blocks.csv")

	if err != nil {

		log.Fatal(err)
	}

	r := csv.NewReader(f)

	for {
		records, err := r.Read()

		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		for value := range records {
			privateBlocks = append(privateBlocks, records[value])
		}
	}

	return privateBlocks
}

func main() {
	// Define the flags
	var help = flag.Bool("help", false, "Show help")

	ip := flag.String("i", "", "Specify a IP address")

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Enable command-line parsing
	flag.Parse()

	i := *ip

	x := &i

	isValid := isValidIp(x)

	list := parsePrivateBlocks()

	if isValid {
		for value := range list {
			subnet_value := list[value]

			clean := strings.Map(func(r rune) rune {
				if unicode.IsGraphic(r) {
					return r
				}
				return -1
			}, subnet_value)

			y := clean

			subnet := &y

			isIn := isInRange(x, subnet)

			// If the IP address is in the range.
			if isIn {
				fmt.Println("Private IP address")
				os.Exit(0)
			}
		}

		fmt.Println("Public IP address")
	} else {
		fmt.Println("Invalid IP address. Use -help for more information.")
	}
}
