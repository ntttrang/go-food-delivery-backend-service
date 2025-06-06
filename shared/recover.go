package shared

import "log"

func Recover() {
	if r := recover(); r != nil {
		log.Println("Recovered:", r)
	}
}
