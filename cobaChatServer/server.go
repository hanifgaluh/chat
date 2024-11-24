package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

var (
	clients = make(map[net.Conn]string)
	mu      sync.Mutex
)

func main() {
	// Input HOST dan PORT
	fmt.Print("Enter server host (default: localhost): ")
	host := readInput("192.168.0.102")
	fmt.Print("Enter server port (default: 33000): ")
	port := readInput("8089")

	address := host + ":" + port
	fmt.Println("Starting server on", address)

	// Mulai server
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server running. Waiting for connections...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		mu.Lock()
		clients[conn] = conn.RemoteAddr().String()
		mu.Unlock()

		fmt.Println("New client connected:", conn.RemoteAddr())
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer func() {
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
		fmt.Println("Client disconnected:", conn.RemoteAddr())
		conn.Close()
		broadcast(fmt.Sprintf("%s has left the chat.", conn.RemoteAddr()), conn)
	}()

	conn.Write([]byte("Welcome to the chat! Type {quit} to exit.\n"))
	broadcast(fmt.Sprintf("%s has joined the chat.", conn.RemoteAddr()), conn)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		if message == "{quit}" {
			break
		}
		broadcast(fmt.Sprintf("%s: %s", conn.RemoteAddr(), message), conn)
	}
}

func broadcast(message string, sender net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	fmt.Println("Broadcasting:", message)
	for conn := range clients {
		if conn != sender {
			conn.Write([]byte(message + "\n"))
		}
	}
}

func readInput(defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	if input == "\n" {
		return defaultValue
	}
	return input[:len(input)-1]
}
