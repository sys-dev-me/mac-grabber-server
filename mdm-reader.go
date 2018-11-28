package main

import "log"
import "fmt"
import "net"
import "os"
import "strings"
import "encoding/json"
import "encoding/binary"
import "bytes"


const (
    CONN_HOST = "10.132.15.213"
    CONN_PORT = "9090"
    CONN_TYPE = "tcp"
)

func  main (){

    // Listen for incoming connections.
    l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)

    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }

    // Close the listener when the application closes.
    defer l.Close()

	mysql_connect()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)

    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        go handleRequest(conn)
    }
}


// Handles incoming requests.
func handleRequest(conn net.Conn) {
 	// Make a buffer to hold incoming data.
 	buf := make([]byte, 1024)

	//var m Authorize

 	// Read the incoming connection into the buffer.
 	reqLen, err := conn.Read(buf)

	log.Println("Received length: ",reqLen)
	log.Println("Data received:", string(buf))

	var pro Authorize
	json.NewDecoder(strings.NewReader(string(buf))).Decode(&pro)
	
	if pro.Request.Type == "auth" {
		answer, _ := respond( pro.Request.Serial )
		answerBody := new(bytes.Buffer)
		err = binary.Write(answerBody, binary.BigEndian, &answer)
		//json.NewEncoder ( answerBody ).Encode ( answer )
		if err != nil {
			log.Println ( "Something went wrong: ", err )
		}
		
		_, err := conn.Write ([]byte(answer))
		if err != nil {
			log.Println("Something went wrong: ", err)
		}
	}

 	if pro.Request.Type == "checkin" {
		var device Request
		json.NewDecoder(strings.NewReader(string(buf))).Decode(&device)

		log.Println( "Request frmo Device: ", device.Request.Serial )
                log.Println( "Hardware type: ", device.Hardware.Model )
                log.Println( "WIFI Interface: ", device.Network.InterfaceName)
		log.Println( "Account: ", device.User.Account )
		log.Println( "Full Name: ", device.User.FullName )
		log.Println( "OS: ", device.OS.OSName )
		log.Println ( "UID:", device.UID )


	}

	log.Println("Request type: ", pro.Request.Type)
	log.Println("Request date: ", pro.Request.Date)
	log.Println("Request frmo Device: ", pro.Request.Serial)


	if err != nil {
   		fmt.Println("Error reading:", err.Error())
 	}
  	// Send a response back to person contacting us.
  	conn.Write([]byte("Message received."))
  	// Close the connection when you're done with it.
 	 conn.Close()
}

