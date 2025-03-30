package main

import (
	"bufio"
	"encoding/json"
	"net"
	"os/exec"
	"testing"
	"time"
)

type TestMessage struct {
	Authentication string `json:"Authentication"`
	MessageType    string `json:"MessageType"`
	Content        CreateUser
}

func startFlatlineService() *exec.Cmd {
	cmd := exec.Command("go", "run", "flatline.go")
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Start()
	time.Sleep(2 * time.Second) // Allow time for server to start
	return cmd
}

func TestServerHeartbeat(t *testing.T) {
	// Start the Flatline service
	cmd := startFlatlineService()
	defer cmd.Process.Kill()

	// Connect to the TCP server
	conn, err := net.Dial("tcp", "127.0.0.1:1313")
	if err != nil {
		t.Fatal("Failed to connect to server:", err)
	}
	defer conn.Close()

	// Create and encode message
	msg := TestMessage{
		Authentication: "I am a little kitty cat",
		MessageType:    "CREATE_INSTITUTION",
		Content:        CreateUser{FirstName: "Aerith", LastName: "Netzer", InstitutionID: 1},
	}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		t.Fatal("Failed to encode JSON:", err)
	}

	// Send message
	writer := bufio.NewWriter(conn)
	writer.Write(jsonMsg)
	writer.WriteByte('\n') // Ensure message is properly terminated
	writer.Flush()

	// Read response
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		t.Fatal("Failed to read response:", err)
	}

	t.Logf("Received: %s", response)
}

func TestCreateUser(t *testing.T) {
	// Start the Flatline service
	cmd := startFlatlineService()
	defer cmd.Process.Kill()

	// Connect to the TCP server
	conn, err := net.Dial("tcp", "127.0.0.1:1313")
	if err != nil {
		t.Fatal("Failed to connect to server:", err)
	}
	defer conn.Close()

	// Create and encode message
	msg := TestMessage{
		Authentication: "TEST_AUTH",
		MessageType:    "CREATE_USER",
	}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		t.Fatal("Failed to encode JSON:", err)
	}

	// Send message
	writer := bufio.NewWriter(conn)
	writer.Write(jsonMsg)
	writer.WriteByte('\n') // Ensure message is properly terminated
	writer.Flush()

	// Read response
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		t.Fatal("Failed to read response:", err)
	}

	t.Logf("Received: %s", response)
}
