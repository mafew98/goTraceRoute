package main

import (
	"TraceRoute/goTrace"
	"flag"
)

// Getting hostname from the user
func main() {
	hostname := flag.String("hostname", "example.com", "Hostname of host to which route is being traced")
	timeout := flag.Int("timeout", 2, "Maximum Number of seconds traceroute spends on each hop before giving up")
	maxHops := flag.Int("maxHops", 64, "Defines the maximum number of hops that will be attempted")
	flag.Parse()

	goTrace.TraceRoute(*hostname, *timeout, *maxHops)
}
