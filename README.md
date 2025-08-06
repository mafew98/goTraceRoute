Project Description:
-------------------
Simple reimplementation of the traceroute network diagnostic tool in Go using UDP, ICMP and raw-socket programming (IPv4, TTL traversal, DNS resolution). Current implementation only uses UDP and IPv4.

Compilation Instructions:
------------------------
1. Build the project using
    sudo go build -o myTraceRoute main.go
2. Run the project using:
    sudo ./myTraceRoute -hostname=dns.google.com -maxHops=32 -timeout=3

** The default values for the command are as follows:
hostname - example.com, maxHops - 64, timeout - 2

Potential Improvements:
-------------------------
The Aim of this project was to complete a simple networking project to learn socket programming in go. The tool can be expanded to support IPv6, TCP and can be designed to run concurrently.

Dependencies:
------------
This project has been tested with go version go1.24.4 darwin/arm64.

Contributing:
------------
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
