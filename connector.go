mdm-reader/mdm-reader.go                                                                            0000664 0002011 0002012 00000005126 13375503717 017153  0                                                                                                    ustar   kravchenko.d                    kravchenko.d                                                                                                                                                                                                           package main

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

                                                                                                                                                                                                                                                                                                                                                                                                                                          mdm-reader/responder.go                                                                             0000644 0000000 0000000 00000002100 13373260675 014126  0                                                                                                    ustar   root                            root                                                                                                                                                                                                                   package main

import "log"
import "time"
import "crypto/rand"
import "math/big"
import "encoding/json"

var possibleChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+,.?/:;{}[]`~")

func generatePassword ( passwordLength int ) string {
	return randomChars ( passwordLength, possibleChars )
} 

func randomChars ( passwordLength int, possibleChars []byte) string {

	newPassword := make ([]byte, passwordLength ) 
	for i := 0; i < passwordLength; i++ {
		j, err := rand.Int(rand.Reader, big.NewInt( int64(len( possibleChars) ) ) )
		if err != nil {
			log.Println( "Something went wrong ", err )
		}
		if int(j.Int64()) == 0 {
			newPassword[i] = possibleChars[0]
			continue
		}
		newPassword[i] = possibleChars[ len ( possibleChars ) % int( j.Int64() ) ]
	}

	return string( newPassword )
}


func respond( serial string ) ([]byte, error) {
	log.Println ( "Send to: ", serial)
	t := time.Now()
	UID := generatePassword( 256 )
	log.Println ( "Generate new UID: ", UID)
	b := AuthAnswer{ serial, t.String(), UID, true}
	return json.Marshal (b)
} 
                                                                                                                                                                                                                                                                                                                                                                                                                                                                mdm-reader/types.go                                                                                 0000644 0000000 0000000 00000002601 13373261203 013263  0                                                                                                    ustar   root                            root                                                                                                                                                                                                                   package main

//general json type of data for mdm checkin request
type Request struct {
     Request     requestType  //timestamp and type of request where send to server
     Hardware    requestHardware // hardware type included serial number device model
     Network     requestWIFI //only wifi,  included name
     User        Account // account type: short and full name
     OS          OSInfo // information aboout device OS, version, build, name
     Versions    Version // version of application and api, static
	UID	string // host key
}

type requestWIFI struct {
    InterfaceName   string
}

type requestHardware struct {
    Model     string
}


// part of checkin request
type Account struct {
    Account  string
    FullName string
}

//part of checking request
type OSInfo struct {
    OSName   string
    OVersion string
    OSBuild  string
}

type requestType struct { 
    Date    string
    Type    string
    Serial      string
}

//part of request
type Version struct {
    Application int 
    API         int 
}

type Authorize struct {
        Request         requestAuth
    Uri string
    isAuthorize bool
}

type requestAuth struct {
    Date    string
    Type    string
    Serial  string
}

type AuthAnswer struct {
        Serial	string	`json:"Serial"`
        Time	string  `json:"Status"`
        UID	string	`json:"UID"`
	ReceivedStatus  bool    `json:"ReceivedStatus"`
}



                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               