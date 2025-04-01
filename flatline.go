package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type Agent struct {
	FirstName     string
	LastName      string
	Scope         string
	InstitutionID int16
}

type Institution struct {
	Name          string
	Address       string
	InstitutionID int16
}

type User struct {
	FirstName     string
	LastName      string
	InstitutionID int16
}

type Message struct {
	Authentication string
	MessageType    string
	Content        User
}

func sendMessage(message Message) {
	conn, err := net.Dial("tcp", "127.0.0.1:1313")
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}
	defer conn.Close()
	jsonMsg, err := json.Marshal(message)
	if err != nil {
		log.Fatal("Failed to encode JSON:", err)
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
		log.Fatal("Failed to read response:", err)
	}

	log.Printf("Received: %s", response)
}

func createUser(message Message) *sql.Rows {
	connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("INSERT INTO users VALUES ($1, $2, $3)",
		message.Content.FirstName, message.Content.LastName)
	return rows
}

func federateInstitution(message Message) *sql.Rows {
	connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("INSERT INTO federatedInstitutions VALUES ($1, $2, $3)",
		message.Content.FirstName, message.Content.LastName)
	return rows
}

func handleMessage(message []byte) {
	var v Message
	err := json.Unmarshal(message, &v)
	if err != nil {
		log.Printf("Failed to unmarshal message: %s", err)
		return
	}
	switch v.MessageType {
	case "CREATE_USER":
		log.Printf("RECEIVED CREATE_USER COMMAND")
		log.Printf("CREATING USER")
		createUser(v)
	case "CREATE_INSTITUTION":
		log.Printf("RECEIVED CREATE_INSTITUTION COMMAND")

	case "FEDERATE_INSTITUTION":
		log.Printf("RECEIVED")
		federateInstitution(v)
	}
	log.Printf("Processed message: %s", message)
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadBytes('\n') // Read until newline
		if err != nil {
			log.Printf("Connection closed: %s", err)
			return
		}
		log.Printf("Message received: %s", message)
		go handleMessage(message)
		fmt.Fprintf(conn, "Message received: %s", message)
	}
}

func main() {
	port := ":1313"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
	defer listener.Close()

	log.Printf("Server listening on %s", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %s", err)
			continue
		}
		go handleConnection(conn)
	}
}
