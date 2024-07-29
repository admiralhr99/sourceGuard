package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

// AllowedIPs is a map of IP addresses that are allowed to access the server.
var AllowedIPs = map[string]bool{
	"192.168.209.116": true, // Example IP, replace with actual allowed IPs
}

func main() {
	listenIP := "192.168.209.124" // Specify the IP address you want to listen on
	listenPort := "50000"         // Specify the port you want to monitor
	listener, err := net.Listen("tcp", listenIP+":"+listenPort)
	if err != nil {
		fmt.Println("Error listening on port:", listenPort, err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Listening on IP:", listenIP, "Port:", listenPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	remoteAddr := clientConn.RemoteAddr().(*net.TCPAddr).IP.String()
	if !AllowedIPs[remoteAddr] {
		fmt.Println("Blocked connection from:", remoteAddr)
		clientConn.Write([]byte("Access denied.\n"))
		return
	}

	targetAddr := "172.31.31.31:2080"
	targetConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		fmt.Println("Error connecting to target:", targetAddr, err)
		return
	}
	defer targetConn.Close()

	fmt.Println("Forwarding connection from:", remoteAddr, "to", targetAddr)

	go io.Copy(targetConn, clientConn)
	io.Copy(clientConn, targetConn)
}
