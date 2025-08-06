package matTraceRoute;

import java.net.InetAddress;
import java.net.SocketException;
import java.net.UnknownHostException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;


public class tracinator {
    private String hostname;
    private InetAddress ipAddress;
    private DatagramSocket socket;
    byte[] buffer = "Im a traceroute packet".getBytes();
    private DatagramPacket tracerPacket;

    public tracinator(String hostname) {
        this.hostname = hostname;
    }

    // get the public ip address from the hostname. Use getByName here.

    private void setHostIPAddress() {
        try {
            ipAddress = InetAddress.getByName(hostname);
        } catch (UnknownHostException uhe) {
            uhe.printStackTrace();
        }
    }

    private void createUDPSocket() throws SocketException{
        this.socket = new DatagramSocket(33434);
    }

    private void createUDPPacket(int TimeToLive) throws SocketException {
        this.tracerPacket = new DatagramPacket(buffer, buffer.length, ipAddress, 33636); // The port here is an arbitrary port and does not matter.
    }

    public void runTrace() throws SocketException {
        // resolve the hostname
        setHostIPAddress();

        // Open a socket and send UDP connections to the IP using TTL as 1.
        createUDPSocket();
        createUDPPacket();

        // Wait to receive ICMP packet back as input on a different socket.
        // Loop to increment TTL to find the number of HOPS.
    } 
}