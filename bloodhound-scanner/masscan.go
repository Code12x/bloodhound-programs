package main

import "github.com/Code12x/go-masscan"

var (
	scanner *masscan.Masscan
)

func setupMasscan() {
	scanner = masscan.New()

	scanner.SetRanges("0.0.0.0-255.255.255.255")
	scanner.SetPorts("25565")
	scanner.SetRate("10000")
	scanner.SetExclude("127.0.0.1")
}
