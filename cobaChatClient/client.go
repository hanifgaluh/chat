package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var (
	clientConn net.Conn
	messages   *widget.Label
	mu         sync.Mutex
)

func main() {
	// Membuat aplikasi Fyne
	chatApp := app.New()
	chatApp.Settings().SetTheme(theme.LightTheme()) // Menggunakan tema terang
	window := chatApp.NewWindow("Chat Client")

	// Menampilkan pesan dengan widget Label
	messages = widget.NewLabel("")
	scroll := container.NewScroll(messages)
	scroll.SetMinSize(fyne.NewSize(400, 300)) // Ukuran minimum area pesan

	// Input untuk pesan
	messageInput := widget.NewEntry()
	messageInput.SetPlaceHolder("Type your message...")

	// Tombol Kirim
	sendButton := widget.NewButton("Send", func() {
		if clientConn != nil {
			text := strings.TrimSpace(messageInput.Text)
			if text == "" {
				return
			}
			sendMessage(text)
			messageInput.SetText("")
		}
	})

	// Layout aplikasi
	window.SetContent(container.NewVBox(
		widget.NewLabelWithStyle("Chat Messages", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		scroll,
		messageInput,
		sendButton,
	))

	// Memulai koneksi ke server
	connectToServer()

	// Jalankan aplikasi
	window.Resize(fyne.NewSize(400, 600))
	window.ShowAndRun()
}

func connectToServer() {
	// Input HOST dan PORT
	fmt.Print("Enter server host (default: 192.168.0.102): ")
	host := readInput("192.168.0.102")
	fmt.Print("Enter server port (default: 8089): ")
	port := readInput("8089")

	address := host + ":" + port
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}

	clientConn = conn
	go receiveMessages()
}

func receiveMessages() {
	scanner := bufio.NewScanner(clientConn)
	for scanner.Scan() {
		message := scanner.Text()
		appendMessage(message)
	}
}

func sendMessage(message string) {
	if message == "{quit}" {
		clientConn.Close()
		os.Exit(0)
	}

	fmt.Fprintln(clientConn, message)
	appendMessage("You: " + message)
}

func appendMessage(message string) {
	mu.Lock()
	defer mu.Unlock()

	current := messages.Text
	messages.SetText(current + "\n" + message)
}

func readInput(defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	if strings.TrimSpace(input) == "" {
		return defaultValue
	}
	return strings.TrimSpace(input)
}
