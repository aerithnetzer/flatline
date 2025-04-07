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
	Content        []byte
}

type SearchMessage struct {
	Authentication string
	MessageType    string
	Content        string
}

func connctToDB() (*sql.DB, error) {
	connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}

func createUser(message Message) []*sql.Rows {
	db, err := connctToDB()
	if err != nil {
		log.Fatal(err)
	}

	type Content struct {
		FirstName     string
		LastName      string
		InstitutionID int16
	}

	var content []Content

	unmarshalerr := json.Unmarshal(message.Content, &content)

	if unmarshalerr != nil {
		fmt.Println("ERROR")
	}
	var success_rows []*sql.Rows
	for _, v := range content {
		rows, qerr := db.Query("INSERT INTO users VALUES ($1, $2, $3)",
			v.FirstName, v.LastName, v.InstitutionID)
		if qerr != nil {
			log.Fatal(qerr)
		}
		success_rows = append(success_rows, rows)
	}

	return success_rows
}

func searchDatabases(message Message) {
}

func requestFederateInstitution(message Message) []*sql.Rows {
	db, err := connctToDB()
	if err != nil {
		log.Fatal(err)
	}

	type Content struct {
		requestingInstitution int16
		myIP                  net.IP
		isSent                bool
	}

	var content []Content

	unmarshalerr := json.Unmarshal(message.Content, &content)

	if unmarshalerr != nil {
		fmt.Println("ERROR")
	}
	var success_rows []*sql.Rows
	for _, v := range content {
		rows, qerr := db.Query("INSERT INTO users VALUES ($1, $2, $3)",
			v.requestingInstitution, v.myIP, v.isSent)
		if qerr != nil {
			log.Fatal(qerr)
		}
		success_rows = append(success_rows, rows)
	}

	return success_rows
}

func acceptFederateinstitution(message Message) []*sql.Rows {
	// Accepts handshake with institution and begins sending federation traffic to the IP

	db, err := connctToDB()
	if err != nil {
		log.Fatal(err)
	}

	type Content struct {
		requestingInstitution int16
		myIP                  net.IP
		isSent                bool
	}

	var content []Content

	unmarshalerr := json.Unmarshal(message.Content, &content)

	if unmarshalerr != nil {
		fmt.Println("ERROR")
	}
	var success_rows []*sql.Rows
	for _, v := range content {
		rows, qerr := db.Query("INSERT INTO users VALUES ($1, $2, $3)",
			v.requestingInstitution, v.myIP, v.isSent)
		if qerr != nil {
			log.Fatal(qerr)
		}
		success_rows = append(success_rows, rows)
	}

	return success_rows
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
	case "REQUEST_FEDERATE_INSTITUTION":
		log.Printf("RECEIVED")
	case "ACCEPT_FEDERATE_INSTITUTION":
		log.Printf("ACCEPTED")

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
