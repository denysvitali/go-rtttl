package main

import (
	rtttl "github.com/denysvitali/go-rtttl/pkg"
	"log"
)

func main(){
	knightRider := "KnightRi:d=4,o=5,b=63:16e6, 32f6, 32e6, 8b6, 16e7, 32f7, 32e7, 8b6, 16e6, 32f6, 32e6, 16b6, 16e7, 4d7, 8p, 4p, 16e6, 32f6, 32e6, 8b6, 16e7, 32f7, 32e7, 8b6, 16e6, 32f6, 32e6, 16b6, 16e7, 4f7, 4p"
	ringtone, err := rtttl.Parse(knightRider)
	if err != nil {
		log.Fatal(err)
	}
	ringtone.Play()
}
