package matTraceRoute;

import org.pcap4j.core.Pcaps;
import org.pcap4j.core.PcapNativeException;

public class CheckPcap {
    public static void main(String[] args) {
        try {
            if (Pcaps.findAllDevs().isEmpty()) {
                System.out.println("No interfaces found. libpcap might be missing.");
            } else {
                System.out.println("libpcap is available!");
            }
        } catch (PcapNativeException e) {
            System.err.println("Error: " + e.getMessage());
            System.err.println("libpcap might not be installed or accessible.");
        }
    }
}