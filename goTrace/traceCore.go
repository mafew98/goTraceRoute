package goTrace

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/ipv4"
)

const packetSize int = 32

var port string = "33444"

type Sender struct {
	target          *net.UDPAddr
	connector       *net.UDPConn
	packetConnector *ipv4.PacketConn
}

type Receiver struct {
	icmpListener net.PacketConn
	buffer       []byte
}

type ResponsePacket struct {
	icmpType  uint8
	icmpCode  uint8
	dstPort   string
	responder net.Addr
}

// Main driver function that initializes entities and begins the packet transfer
func TraceRoute(hostname string, timeout int, maxHops int) {
	var sender Sender
	var receiver Receiver

	initializeSender(hostname, &sender, maxHops)
	initializeReceiver(&receiver)

	// Set TTL in connection and serve data
	var hop int
	for hop = 1; hop <= maxHops; hop++ {
		traceByHop(&sender, hop)
		hopStart := time.Now()
		retval := listenForReply(&receiver, hop)
		rtt := time.Since(hopStart)
		if retval != 2 {
			fmt.Printf("    %.3f ms\n", float64(rtt.Nanoseconds())/1e6)
			if retval == 1 {
				break
			}
		}
	}

	if sender.connector != nil {
		sender.connector.Close()
	}
	if receiver.icmpListener != nil {
		receiver.icmpListener.Close()
	}
}

func initializeSender(hostname string, sender *Sender, maxHops int) {
	// Acquiring the first IPv4 address for the given hostname
	ips, err := net.LookupIP(hostname)
	if err != nil {
		fmt.Println("Unable to lookup IPs for the given hostname " + hostname)
	}

	// Setting to the first IPv4 address associated to the hostname
	var targetIP string
	for _, ip := range ips {
		if ip.To4() != nil {
			targetIP = ip.String()
			break
		}
	}
	if targetIP == "" {
		log.Fatal("Unable to resolve any IPv4 addresses for the given hostname ", hostname)
	}

	sender.target, err = net.ResolveUDPAddr("udp", net.JoinHostPort(targetIP, port))
	if err != nil {
		log.Fatal(err)
	}

	// Print connections parameters acquired
	fmt.Printf("traceroute to %v (%v), %d hops max, %d byte packets\n", hostname, sender.target.IP, maxHops, packetSize)

	// Establish connection using parameters
	// Setting up socket for UDP connection
	// sender.connector, err = net.DialUDP("udp", nil, sender.target)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	UDPConn, err := net.ListenPacket("udp4", ":0") // Use ephemeral source port since sending port does not matter
	if err != nil {
		log.Fatal(err)
	}
	sender.packetConnector = ipv4.NewPacketConn(UDPConn)
}

func initializeReceiver(receiver *Receiver) {
	var err error
	// Create ICMP Reply listener socket
	receiver.icmpListener, err = net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatal("Unable to create ICMP Listener\n", err)
		os.Exit(126)
	}

	// Allocate Buffer
	receiver.buffer = make([]byte, 1500)
}

func traceByHop(sender *Sender, hops int) {
	// Setting the TTL for the packet
	var err error = sender.packetConnector.SetTTL(int(hops))
	if err != nil {
		log.Fatal(err)
		os.Exit(126)
	}

	// Send the packet
	message := []byte("Mat custom trace route")
	_, err = sender.packetConnector.WriteTo(message, nil, sender.target)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println("Traceroute Packet Sent!")
}

/*
Returns false if traceroute hasn't reached the target and true if it has reached the final destination.
*/
func listenForReply(receiver *Receiver, step int) int {
	for {
		receiver.icmpListener.SetReadDeadline(time.Now().Add(time.Second))
		_, addr, err := receiver.icmpListener.ReadFrom(receiver.buffer)
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				// timeout is expected
				fmt.Printf("%d  * * *\n", step)
				return 2
			} else {
				// non-timeout errors are fatal
				log.Fatalln(err)
				os.Exit(1)
			}
		}

		// fmt.Println("received", n, "bytes from", addr)
		// Got some ICMP response. Time to validate it.
		responsePacket, err := parseResponse(receiver)
		if err != nil {
			// ignoring incorrect packet
			continue
		}
		responsePacket.responder = addr

		if isValidResponse(responsePacket) {
			// fmt.Println("Got a valid response")
			if handleResponse(step, responsePacket) {
				return 1
			} else {
				return 0
			}
		} else {
			continue
		}
	}
}

/*
Returns parsed response packet and error
If there is no error, indicates that a valid reply was obtained.
*/
func parseResponse(receiver *Receiver) (response *ResponsePacket, err error) {
	if len(receiver.buffer) < 36 {
		err = fmt.Errorf("packet too short to be valid ICMP with embedded IP+UDP")
		return nil, err
	}
	response = &ResponsePacket{}
	response.icmpType = receiver.buffer[0]
	response.icmpCode = receiver.buffer[1]

	// Skip ICMP header → embedded IP header starts at byte 8
	ipHeaderStart := 8
	udpHeaderStart := ipHeaderStart + 20

	// UDP header is 8 bytes: src port (2), dst port (2), length (2), checksum (2)
	response.dstPort = strconv.Itoa(int(binary.BigEndian.Uint16(receiver.buffer[udpHeaderStart+2 : udpHeaderStart+4])))

	return response, nil
}

/*
Since the ICMP listener receives all icmp packets, we need to ignore the ones that are not related to the traceroute ping.
This is done by comparing the sending port number with that of the received packet.
*/
func isValidResponse(response *ResponsePacket) bool {
	// fmt.Println("Destination port is ", response.dstPort)
	return (response.dstPort == port)
}

// Handles a valid packet response to traceroute and assigns the appropriate actions.
func handleResponse(step int, response *ResponsePacket) bool {
	hostArr, err := net.LookupAddr(response.responder.String())
	var host string
	if err != nil || len(hostArr) == 0 {
		// Unable to find the hostname tied to the received IP address.
		host = response.responder.String()
	} else {
		host = hostArr[0]
	}
	fmt.Printf("%d  %v (%v)", step, host, response.responder)

	if response.icmpType == 11 && response.icmpCode == 0 {
		// fmt.Println("TLL Expired detected")
		// TTL expired in transit.
		return false
	} else if response.icmpCode == 3 && response.icmpType == 3 {
		// Reached final destination (destination unreachable (port))
		return true
	}
	return false
}
