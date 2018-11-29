package main

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
