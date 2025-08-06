package main

import (
	"TraceRoute/goTrace"
	"flag"
)

// Getting hostname from the user
func main() {
	hostname := flag.String("hostname", "example.com", "Hostname of host to which route is being traced")
	flag.Parse()

	goTrace.TraceRoute(*hostname)
}
