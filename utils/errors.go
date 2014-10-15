package utils

import (
	"log"
)

func Check(err error) {
	if err != nil {
		log.Printf("Warning: an error occured")
		log.Printf("%v", err)
	}
}
