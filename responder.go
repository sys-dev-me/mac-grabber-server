package main

import "log"
import "time"
import "crypto/rand"
import "math/big"
import "encoding/json"

type sessionKey struct {
	uid []byte
}


var possibleChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+,.?/:;{}[]`~")

func generatePassword ( passwordLength int ) string {
	pw := randomChars ( passwordLength, possibleChars )
	return string(pw)
} 

func randomChars ( passwordLength int, possibleChars []byte) []byte {

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

	return newPassword
}

var tmp []string

func respond( serial string ) ([]byte, error) {
	if tmp == nil {
		tmp = make([]string, 0)
	}

	log.Println ( "Send to: ", serial)
	t := time.Now()
		UID := generatePassword( 256 )
		tmp = append( tmp, UID)

	log.Println ( "Generate new UID: ", UID )
	b := AuthAnswer{ serial, t.String(), UID, true}


	log.Println ( "pw stotrage: " )
	for i :=0 ; i < len (tmp ); i ++ {
		log.Println ( i, " => ", tmp[i] )
	}

	return json.Marshal (b)
} 
