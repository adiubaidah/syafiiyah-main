package main

import (
	"time"

	"github.com/adiubaidah/rfid-syafiiyah/initiator"
)

func main() {

	time.Local = time.UTC
	initiator.Init()
}
